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

	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/api/paths"
	"github.com/cjslep/dharma/internal/features"
	"github.com/cjslep/dharma/internal/render"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

func (e *ESIAuth) getAuth(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
	fs, err := e.C.Features.GetEnabledFeatures(langs...)
	if err != nil {
		e.C.MustRenderError(w, r, errors.Wrap(err, "could not fetch enabled features"), langs...)
		return
	}
	l := features.List(fs)

	u := paths.GetPostESIAuthPath(langs[0], l)
	rc := api.From(r.Context())
	v := render.NewHTMLView(
		w,
		http.StatusOK,
		"esiauth/auth",
		rc,
		map[string]interface{}{
			"scopes":   l.Scopes(),
			"explain":  l.ScopeExplanations(),
			"authPath": u.String(),
		},
		langs...)
	e.C.MustRender(v)
}
