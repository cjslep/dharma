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
	"context"
	"time"

	"github.com/rs/zerolog"
)

// Periodic functions are ran intermittently, and errors returned from them
// result in logged error messages.
type PeriodicFn func(context.Context) error

type periodic struct {
	D     time.Duration
	F     PeriodicFn
	T     *time.Timer
	M     *Messenger
	L     *zerolog.Logger
	First bool
}

func (m *Messenger) Periodically(d time.Duration, f PeriodicFn, l *zerolog.Logger) {
	p := &periodic{
		D:     d,
		F:     f,
		T:     time.NewTimer(d),
		M:     m,
		L:     l,
		First: false,
	}
	go p.run()
}

func (m *Messenger) NowAndPeriodically(d time.Duration, f PeriodicFn, l *zerolog.Logger) {
	p := &periodic{
		D:     d,
		F:     f,
		T:     time.NewTimer(d),
		M:     m,
		L:     l,
		First: true,
	}
	go p.run()
}

func (p *periodic) run() {
	if p.First {
		err := p.F(p.M.parentCtx)
		if err != nil {
			p.L.Error().Stack().Err(err).Msg("")
		}
	}
	for {
		select {
		case <-p.T.C:
			err := p.F(p.M.parentCtx)
			if err != nil {
				p.L.Error().Stack().Err(err).Msg("")
			}
			p.T.Reset(p.D)
		case <-p.M.parentCtx.Done():
			if !p.T.Stop() {
				<-p.T.C
			}
			return
		}
	}
}
