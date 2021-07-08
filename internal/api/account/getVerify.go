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
	"fmt"
	"net/http"
	"net/url"

	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/render"
	"github.com/go-fed/apcore/app"
	"golang.org/x/text/language"
)

func getVerifyURL(lang language.Tag, showThanksForRegistering bool) *url.URL {
	u := &url.URL{}
	u.Path = fmt.Sprintf("/%s/verify", lang)
	v := url.Values{}
	if showThanksForRegistering {
		v.Add("ty", "true")
	}
	u.RawQuery = v.Encode()
	return u
}

func (a *Account) getVerify(w http.ResponseWriter, r *http.Request, k app.Session, langs []language.Tag) {
	showTY := r.URL.Query().Get("ty")
	rc := api.From(r.Context())
	v := render.NewHTMLView(
		w,
		http.StatusOK,
		"account/verify",
		rc,
		map[string]interface{}{
			"showty": showTY,
		},
		langs...)
	a.C.MustRender(v)
}
