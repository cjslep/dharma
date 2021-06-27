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

package util

import (
	"context"
	"sync"
	"time"
)

var _ context.Context = &mergedContext{}

// mergedContext is able to function as a context.Context, but with 2 parents.
//
// See https://github.com/golang/go/issues/36503
type mergedContext struct {
	parent0 context.Context
	parent1 context.Context
	done    chan struct{}
	// For coordinating the error
	mu  sync.RWMutex
	err error
}

func Merge(p0, p1 context.Context) context.Context {
	mc := &mergedContext{
		parent0: p0,
		parent1: p1,
		done:    make(chan struct{}),
	}
	go mc.run()
	return mc
}

func (m *mergedContext) Deadline() (deadline time.Time, ok bool) {
	d0, ok0 := m.parent0.Deadline()
	d1, ok1 := m.parent1.Deadline()
	dl := d0
	if d0.IsZero() || (!d1.IsZero() && d1.Before(d0)) {
		dl = d1
	}
	return dl, ok0 || ok1
}

func (m *mergedContext) Done() <-chan struct{} {
	return m.done
}

func (m *mergedContext) run() {
	var te error
	select {
	case <-m.parent0.Done():
		te = m.parent0.Err()
	case <-m.parent1.Done():
		te = m.parent1.Err()
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.err = te
	close(m.done)
}

func (m *mergedContext) Err() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.err
}

func (m *mergedContext) Value(key interface{}) interface{} {
	v0 := m.parent0.Value(key)
	if v0 == nil {
		return m.parent1.Value(key)
	}
	return v0
}
