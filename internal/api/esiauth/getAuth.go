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
	"github.com/cjslep/dharma/internal/util"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

func (e *ESIAuth) getAuth(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
	ctx := e.C.F.Context(r)
	var fs []features.Feature
	if e.C.State.RequiresCorpToBeManaged() {
		var err error
		fs, err = e.C.Features.GetAdminCEOInitialFeatures(ctx, langs...)
		if err != nil {
			e.C.MustRenderError(w, r, errors.Wrap(err, "could not fetch features for administration/CEO"), langs...)
			return
		}
	} else {
		var err error
		fs, err = e.C.Features.GetEnabledFeatures(ctx, langs...)
		if err != nil {
			e.C.MustRenderError(w, r, errors.Wrap(err, "could not fetch enabled features"), langs...)
			return
		}
	}
	l := features.List(fs)

	isRescopeText := r.URL.Query().Get(paths.RescopeQueryParam)
	isRescope := len(isRescopeText) > 0

	u := paths.GetPostESIAuthPath(util.GetPreferredLanguage(langs), l)
	rc := api.From(ctx)
	v := render.NewHTMLView(
		w,
		http.StatusOK,
		"esiauth/auth",
		rc,
		map[string]interface{}{
			"scopes":    l.Scopes(),
			"explain":   l.ScopeExplanations(),
			"authPath":  u.String(),
			"isRescope": isRescope,
		},
		langs...)
	e.C.MustRender(v)
}
