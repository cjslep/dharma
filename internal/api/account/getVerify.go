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
	"golang.org/x/text/language"
)

func (a *Account) getVerify(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
	token := r.URL.Query().Get(paths.TokenQueryParam)
	showTY := r.URL.Query().Get(paths.VerifyTYParam)
	success := r.URL.Query().Get(paths.VerifySuccessParam)
	rc := api.From(r.Context())
	v := render.NewHTMLView(
		w,
		http.StatusOK,
		"account/verify",
		rc,
		map[string]interface{}{
			"token":    token,
			"showty":   showTY,
			"verified": success,
		},
		langs...)
	a.C.MustRender(v)
}
