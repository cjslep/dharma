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
	"github.com/pkg/errors"
)

func (e *ESIAuth) getCallback(w http.ResponseWriter, r *http.Request) {
	k, err := e.C.F.Session(r)
	if err != nil {
		e.C.ErrorHandler(w, r, errors.New("could not retrieve session"))
		return
	}

	// Enforce state integrity
	state := r.URL.Query().Get("state")
	myState := sessions.GetESIOAuth2State(k)
	sessions.ClearESIOAuth2State(k)
	if myState == "" || state != myState {
		if err := k.Save(r, w); err != nil {
			e.C.ErrorHandler(w, r, errors.Wrap(err, "could not save session"))
			return
		}
		e.C.ErrorHandler(w, r, errors.New("oauth2 state mismatch"))
		return
	}

	// Exchange the short-lived code for long-term authorization.
	code := r.URL.Query().Get("code")
	jwt, err := e.C.OAC.GetAuthorization(code)
	if err != nil {
		e.C.ErrorHandler(w, r, errors.Wrap(err, "could not get jwt"))
		return
	}

	// Verify the authenticity of the authorization.
	ek, err := e.C.DB.GetEvePublicKeys()
	if err != nil {
		e.C.ErrorHandler(w, r, errors.Wrap(err, "could not retrieve EVE public keys"))
		return
	}
	jwtk := ek.JWTKey()
	if jwtk == nil {
		e.C.ErrorHandler(w, r, errors.New("could not find EVE public jwt key"))
		return
	}
	claims, err := jwtk.ValidateToken([]byte(jwt.AccessToken))
	if err != nil {
		e.C.ErrorHandler(w, r, errors.Wrap(err, "could not validate jwt"))
		return
	}

	// Construct our internal representation of a validated token, and
	// store it.
	tokens, err := esi.NewTokens(jwt, claims)
	if err != nil {
		e.C.ErrorHandler(w, r, errors.Wrap(err, "could not construct tokens"))
		return
	}
	err = e.C.DB.SetEveTokens(tokens)
	if err != nil {
		e.C.ErrorHandler(w, r, errors.Wrap(err, "could not store tokens"))
		return
	}
	// TODO: Launch periodic jobs to refresh expiring access tokens.

	// Finally, write the response to the awaiting connection.
	if err := k.Save(r, w); err != nil {
		e.C.ErrorHandler(w, r, errors.Wrap(err, "could not save session"))
		return
	}
	http.Redirect(w, r /*TODO*/, "", http.StatusFound)
}