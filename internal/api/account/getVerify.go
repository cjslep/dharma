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

package account

import (
	"net/http"

	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/api/paths"
	"github.com/cjslep/dharma/internal/render"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

func (a *Account) getVerify(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
	token := r.URL.Query().Get(paths.TokenQueryParam)
	if token == "" {
		// If no token, then display page directing user to verify
		// their email.
		showTY := r.URL.Query().Get("ty")
		rc := api.From(r.Context())
		v := render.NewHTMLView(
			w,
			http.StatusOK,
			"account/verify",
			rc,
			map[string]interface{}{
				"verified": false,
				"showty":   showTY,
			},
			langs...)
		a.C.MustRender(v)
		return
	}

	// If token is present, instead attempt to verify.
	err := a.C.Users.MarkUserVerified(a.C.F.Context(r), token)
	if err != nil {
		a.C.MustRenderError(w, r, errors.Wrap(err, "could not verify user"), langs...)
		return
	}
	v := render.NewHTMLView(
		w,
		http.StatusOK,
		"account/verify",
		api.From(r.Context()),
		map[string]interface{}{
			"verified": true,
			"showty":   false,
		},
		langs...)
	a.C.MustRender(v)
}
