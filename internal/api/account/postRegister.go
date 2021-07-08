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
	"github.com/cjslep/dharma/internal/render"
	"github.com/go-fed/apcore/app"
	"github.com/mholt/binding"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type registerRequest struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

func (r *registerRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&r.Username: binding.Field{
			Form:     "username",
			Required: true,
		},
		&r.Email: binding.Field{
			Form:     "email",
			Required: true,
		},
		&r.Password: binding.Field{
			Form:     "password",
			Required: true,
		},
		&r.ConfirmPassword: binding.Field{
			Form:     "confirm_password",
			Required: true,
		},
	}
}

// TODO: Rate limit, due to username/email taken errors

func (a *Account) postRegister(w http.ResponseWriter, r *http.Request, k app.Session, langs []language.Tag) {
	rc := api.From(r.Context())
	rr := &registerRequest{}
	errs := binding.Bind(r, rr)
	if errs.Len() > 0 {
		v := render.NewBadRequestView(w, rc, langs...)
		a.C.MustRender(v)
		return
	}

	if rr.Password != rr.ConfirmPassword {
		u := getRegisterURLPasswordsDontMatch(r, rr.Username, rr.Email)
		http.Redirect(w, r, u.String(), http.StatusFound)
		return
	}

	err := a.C.Users.CreateUser(a.C.F.Context(r), rr.Username, rr.Email, rr.Password)
	if err != nil {
		if a.C.F.IsNotUniqueEmail(err) {
			u := getRegisterURLEmailNotUnique(r, rr.Username, rr.Email)
			http.Redirect(w, r, u.String(), http.StatusFound)
			return
		} else if a.C.F.IsNotUniqueUsername(err) {
			u := getRegisterURLUsernameNotUnique(r, rr.Username, rr.Email)
			http.Redirect(w, r, u.String(), http.StatusFound)
			return
		} else {
			a.C.MustRenderError(w, r, errors.Wrap(err, "could not create user"), langs...)
			return
		}
	}

	u := getVerifyURL(langs[0], true)
	http.Redirect(w, r, u.String(), http.StatusFound)
}
