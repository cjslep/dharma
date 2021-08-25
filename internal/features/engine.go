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
	"sort"

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

func (e *Engine) ValidateFeatureIDs(ids []string) error {
	for _, id := range ids {
		if !allFeatureIDs[id] {
			return errors.Errorf("invalid feature id: %s", id)
		}
	}
	return nil
}

func (e *Engine) DiffScopes(currentIDs, enableIDs, disableIDs []string, langs ...language.Tag) (added, removed []Scope, err error) {
	// 0. Validate inputs are well-formed
	em := make(map[string]bool, len(enableIDs))
	for _, id := range enableIDs {
		em[id] = true
	}
	dm := make(map[string]bool, len(disableIDs))
	for _, id := range disableIDs {
		dm[id] = true
	}
	// Allowed:
	//   - Enabling already-enabled
	//   - Disabling already-disabled
	// Disallowed:
	//   - Both enabling and disabling
	for _, id := range enableIDs {
		if dm[id] {
			err = errors.Errorf("diffing enabling & disabling same feature: %s", id)
			return
		}
	}
	for _, id := range disableIDs {
		if em[id] {
			err = errors.Errorf("diffing enabling & disabling same feature: %s", id)
			return
		}
	}
	// 1. Get current feature set
	var current, next []Feature
	current, err = e.GetFeatures(currentIDs, langs...)
	if err != nil {
		return
	}
	// 2. Get next feature set
	cm := make(map[string]bool, len(currentIDs))
	for _, id := range currentIDs {
		cm[id] = true
	}
	for _, id := range enableIDs {
		cm[id] = true
	}
	for _, id := range disableIDs {
		cm[id] = false
	}
	nextIDs := make([]string, 0, len(cm))
	for id, enabled := range cm {
		if enabled {
			nextIDs = append(nextIDs, id)
		}
	}
	next, err = e.GetFeatures(nextIDs, langs...)
	if err != nil {
		return
	}
	// 3. Diff scopes
	cs := List(current).Scopes()
	ns := List(next).Scopes()
	csm := make(map[Scope]bool, len(cs))
	nsm := make(map[Scope]bool, len(ns))
	for _, v := range cs {
		csm[v] = true
	}
	for _, v := range ns {
		nsm[v] = true
	}
	for k := range nsm {
		if !csm[k] {
			added = append(added, k)
		}
	}
	for k := range csm {
		if !nsm[k] {
			removed = append(removed, k)
		}
	}
	// 4. Sort for deterministic results
	sort.Sort(Scopes(added))
	sort.Sort(Scopes(removed))
	return
}
