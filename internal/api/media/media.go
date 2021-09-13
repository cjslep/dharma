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

package media

import (
	"net/http"

	"github.com/cjslep/dharma/internal/api"
	"github.com/go-fed/apcore/app"
)

type Media struct {
	C *api.Context
}

func (m *Media) Route(r app.Router) {
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/media/{id}",
		api.MustHaveSession(m.C, m.getMedia))
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/media/portraits/{id}",
		http.HandlerFunc(m.getCharacterPortrait))
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/media/corporations/{id}",
		http.HandlerFunc(m.getCorporationIcon))
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/media/alliances/{id}",
		http.HandlerFunc(m.getAllianceIcon))
	r.NewRoute().Methods("POST").WebOnlyHandler(
		"/media",
		api.MustHaveSession(m.C, m.postMedia))
}
