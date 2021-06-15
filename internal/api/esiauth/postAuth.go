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
	"net/http"

	"github.com/cjslep/dharma/internal/sessions"
)

func (e *ESIAuth) postAuth(w http.ResponseWriter, r *http.Request) {
	k, err := e.C.F.Session(r)
	if err != nil {
		// TODO
		return
	}
	var state string    // TODO
	var scopes []string // TODO
	sessions.SetESIOAuth2State(k, state)
	u := e.C.OAC.GetURL(state, scopes)
	if err := k.Save(r, w); err != nil {
		// TODO
		return
	}
	http.Redirect(w, r, u.String(), http.StatusFound)
}
