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
	"net/url"

	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/render"
	"golang.org/x/text/language"
)

func getRegisterURLPasswordsDontMatch(r *http.Request, username, email string) *url.URL {
	return getRegisterURL(r, username, email, "passwordsDoNotMatch")
}

func getRegisterURLUsernameNotUnique(r *http.Request, username, email string) *url.URL {
	return getRegisterURL(r, username, email, "usernameNotUnique")
}

func getRegisterURLEmailNotUnique(r *http.Request, username, email string) *url.URL {
	return getRegisterURL(r, username, email, "emailNotUnique")
}

func getRegisterURL(r *http.Request, username, email, err string) *url.URL {
	u := &url.URL{}
	u.Path = r.URL.Path
	v := url.Values{}
	v.Add("u", username)
	v.Add("e", email)
	v.Add("err", err)
	u.RawQuery = v.Encode()
	return u
}

func (a *Account) getRegister(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
	username := r.URL.Query().Get("u")
	email := r.URL.Query().Get("e")
	err := r.URL.Query().Get("err")
	rc := api.From(r.Context())
	v := render.NewHTMLView(
		w,
		http.StatusOK,
		"account/register",
		rc,
		map[string]interface{}{
			"username": username,
			"email":    email,
			"err":      err,
		},
		langs...)
	a.C.MustRender(v)
}
