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

package esiauth

import (
	"github.com/cjslep/dharma/internal/api"
	"github.com/go-fed/apcore/app"
)

const (
	Callback = "/esi/callback"
)

type ESIAuth struct {
	C *api.Context
}

func (e *ESIAuth) Route(r app.Router) {
	r.NewRoute().Methods("GET").WebOnlyHandler("/esi/auth", api.MustHaveLanguageCode(e.C, e.getAuth))
	r.NewRoute().Methods("POST").WebOnlyHandler("/esi/auth", api.MustHaveSessionAndLanguageCode(e.C, e.postAuth))
	r.NewRoute().Methods("GET").WebOnlyHandler(Callback, api.MustHaveSession(e.C, e.getCallback))
}
