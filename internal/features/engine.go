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
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type Engine struct {
	b *i18n.Bundle
}

func New(b *i18n.Bundle) *Engine {
	e := &Engine{
		b: b,
	}
	return e
}

func (e *Engine) localized(langs ...language.Tag) ([]Feature, error) {
	l := make([]string, len(langs))
	for i := range langs {
		l[i] = langs[i].String()
	}

	return allLocalizedFeatures(e.b, l...)
}

func (e *Engine) GetRequiredFeatures(langs ...language.Tag) ([]Feature, error) {
	f, err := e.localized(langs...)
	if err != nil {
		return nil, err
	}
	var out []Feature
	for _, t := range f {
		if t.Required {
			out = append(out, t)
		}
	}
	return out, nil
}

func (e *Engine) GetFeatures(ids []string, langs ...language.Tag) ([]Feature, error) {
	m := make(map[string]bool, len(ids))
	for _, v := range ids {
		m[v] = true
	}

	f, err := e.localized(langs...)
	if err != nil {
		return nil, err
	}
	var out []Feature
	for _, t := range f {
		if m[t.ID] {
			out = append(out, t)
		}
	}
	return out, nil
}

func (e *Engine) GetAllFeatures(langs ...language.Tag) ([]Feature, error) {
	return e.localized(langs...)
}

func (e *Engine) ValidateFeatureID(id string) error {
	ok := allFeatureIDs[id]
	if !ok {
		return errors.Errorf("invalid feature id: %s", id)
	}
	return nil
}
