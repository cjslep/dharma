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
	"github.com/cjslep/dharma/internal/features"
	"github.com/cjslep/dharma/internal/sessions"
	"github.com/go-fed/apcore/app"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

func (e *ESIAuth) postAuth(w http.ResponseWriter, r *http.Request, k app.Session, langs []language.Tag) {
	f, ok := r.URL.Query()["features"]
	if !ok {
		e.C.MustRenderError(w, r, errors.New("could not get any features selected"), langs...)
		return
	}
	fl, err := e.C.Features.GetFeatures(f)
	if err != nil {
		e.C.MustRenderError(w, r, errors.Wrap(err, "could not obtain features"), langs...)
		return
	}
	scopes := features.List(fl).Scopes()

	state, err := esi.Random(64)
	if err != nil {
		e.C.MustRenderError(w, r, errors.Wrap(err, "could not generate state for oauth2"), langs...)
		return
	}
	sessions.SetESIOAuth2State(k, state)
	u := e.C.OAC.GetURL(state, scopes)
	if err := k.Save(r, w); err != nil {
		e.C.MustRenderError(w, r, errors.Wrap(err, "could not save session"), langs...)
		return
	}
	http.Redirect(w, r, u.String(), http.StatusFound)
}
