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

package site

import (
	"github.com/cjslep/dharma/internal/async"
	"github.com/cjslep/dharma/internal/db"
	"github.com/go-fed/apcore/app"
)

type Site struct {
	db *db.DB
	m  *async.Messenger
}

func New(db *db.DB, m *async.Messenger) *Site {
	return &Site{
		db: db,
		m:  m,
	}
}

func (s *Site) Route(r app.Router) {
	// TODO
	r.Methods("GET").WebOnlyHandlerFunc("/", s.getHome)
	r.Methods("GET").WebOnlyHandlerFunc("/about", s.getAbout)
}
