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

package features

import (
	"github.com/pkg/errors"
)

type Feature struct {
	Name        string
	Description string
	Scopes      []ScopeExplanation
	Required    bool
}

type ScopeExplanation struct {
	Scope       string
	Explanation string
}

type Engine struct {
	f     []Feature
	names map[string]int
}

func New(f []Feature) *Engine {
	e := &Engine{
		f:     f,
		names: make(map[string]int, 0),
	}
	for i, t := range f {
		e.names[t.Name] = i
	}
	return e
}

func (e *Engine) Convert(features []string) (scopes []string, err error) {
	// First, build map of requested scopes
	m := make(map[string]bool, 1)
	applied := make(map[string]bool, 1)
	for _, f := range features {
		idx, ok := e.names[f]
		if !ok {
			err = errors.New("no feature with name: " + f)
			return
		}
		for _, s := range e.f[idx].Scopes {
			m[s.Scope] = true
		}
		applied[f] = true
	}

	// Next, ensure all required features were applied
	for _, f := range e.f {
		_, ok := applied[f.Name]
		if f.Required && !ok {
			err = errors.New("required feature was not requested: " + f.Name)
			return
		}
	}

	// Finally, convert the map to a list
	scopes = make([]string, 0, len(m))
	for f, _ := range m {
		scopes = append(scopes, f)
	}
	return
}

func (e *Engine) GetFeatures() []Feature {
	return e.f
}
