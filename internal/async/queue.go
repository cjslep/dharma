// dharma is a supplementary corporation community tool for Eve Online.
// Copyright (C) 2021 Cory Slep
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package async

import (
	"errors"
	"sync"
)

// CallbackFn is a callback function signature.
//
// Since Go does not (yet) have generics, it allows encapsulating references to
// specific types in the caller.
type CallbackFn func() error

// DoFn is the basic unit of async execution.
type DoFn func() CallbackFn

// message is a generic message to handle asynchronously.
type message struct {
	// Do is the unit of work to do asynchronously.
	Do DoFn
	// Out allows passing back a callback once async processing is done.
	Out chan<- CallbackFn
}

// Messenger is held by clients who wish to execute something asynchronously.
type Messenger struct {
	closed bool
	in     chan<- message
	done   chan bool
}

// DoAsync executes the DoFn asynchronously, returning a channel which will have
// the callback function provided upon completion.
//
// If the Messenger is already closed, a callback function returning an error is
// given.
func (m *Messenger) DoAsync(f DoFn) <-chan CallbackFn {
	cb := make(chan CallbackFn, 1)
	if m.isClosed() {
		cb <- func() error {
			return errors.New("Queue channel already closed")
		}
		return cb
	}
	m.in <- message{
		Do:  f,
		Out: cb,
	}
	return cb
}

// DoBlocking executes the function in the Queue queue, blocks until a callback
// is received, then executes the callback and returns any errors.
func (m *Messenger) DoBlocking(f DoFn) error {
	cb := <-m.DoAsync(f)
	return cb()
}

func (m *Messenger) isClosed() bool {
	_, ok := <-m.done
	return !ok
}

// Queue gracefully facilitates fan-in message passing.
type Queue struct {
	stopOnce  sync.Once
	startOnce sync.Once
	c         chan message
	m         []*Messenger
	done      chan bool
	ackDone   chan bool
}

func NewQueue() *Queue {
	return &Queue{
		c:       make(chan message),
		m:       make([]*Messenger, 0),
		done:    make(chan bool),
		ackDone: make(chan bool),
	}
}

// Start begins the asynchronous processing.
//
// It must be called at the beginning of the Queue lifetime. Starting
// after a call to Stop is not supported as all Messengers become invalidated.
func (a *Queue) Start() error {
	if a.closed() {
		return errors.New("Queue does not support Start after Stop")
	}
	a.startOnce.Do(func() {
		go func() {
		GoLoop:
			for {
				select {
				case do := <-a.c:
					cb := do.Do()
					do.Out <- cb
				case <-a.done:
					break GoLoop
				}
			}
			for do := range a.c {
				cb := do.Do()
				do.Out <- cb
			}
			close(a.ackDone)
		}()
	})
	return nil
}

// Stop ends asynchronous processing gracefully, returning after all processing
// is completed.
//
// Invalidates the Messengers produced by this Queue.
func (a *Queue) Stop() {
	if a.closed() {
		return
	}
	a.stopOnce.Do(func() {
		close(a.done)
		<-a.ackDone
	})
}

// Messenger obtains a new handle for clients to add work to the async queue.
//
// If the Queue pool is already closed, nil is returned.
func (a *Queue) Messenger() *Messenger {
	if a.closed() {
		return nil
	}
	return &Messenger{
		in:   a.c,
		done: a.done,
	}
}

func (a *Queue) closed() bool {
	select {
	case _, ok := <-a.done:
		return !ok
	default:
		return false
	}
}
