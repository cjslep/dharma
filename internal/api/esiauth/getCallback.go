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

	"github.com/cjslep/dharma/esi"
	"github.com/cjslep/dharma/internal/sessions"
)

func (e *ESIAuth) getCallback(w http.ResponseWriter, r *http.Request) {
	k, err := e.f.Session(r)
	if err != nil {
		// TODO
		return
	}

	// Enforce state integrity
	state := r.URL.Query().Get("state")
	myState := sessions.GetESIOAuth2State(k)
	sessions.ClearESIOAuth2State(k)
	if myState == "" || state != myState {
		if err := k.Save(r, w); err != nil {
			// TODO
			return
		}
		// TODO
		return
	}

	// Exchange the short-lived code for long-term authorization.
	code := r.URL.Query().Get("code")
	jwt, err := e.oac.GetAuthorization(code)
	if err != nil {
		// TODO
		return
	}

	// Verify the authenticity of the authorization.
	ek, err := e.db.GetEvePublicKeys()
	if err != nil {
		// TODO
		return
	}
	jwtk := ek.JWTKey()
	if jwtk == nil {
		// TODO
		return
	}
	claims, err := jwtk.ValidateToken([]byte(jwt.AccessToken))
	if err != nil {
		// TODO
		return
	}

	// Construct our internal representation of a validated token, and
	// store it.
	tokens, err := esi.NewTokens(jwt, claims)
	if err != nil {
		// TODO
		return
	}
	err = e.db.SetEveTokens(tokens)
	if err != nil {
		// TODO
		return
	}
	// TODO: Launch periodic jobs to refresh expiring access tokens.

	// Finally, write the response to the awaiting connection.
	if err := k.Save(r, w); err != nil {
		// TODO
	}
	http.Redirect(w, r /*TODO*/, "", http.StatusFound)
}
