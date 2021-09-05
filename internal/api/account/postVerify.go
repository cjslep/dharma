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
	"context"
	"net/http"

	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/api/paths"
	"github.com/cjslep/dharma/internal/async"
	"github.com/cjslep/dharma/internal/render"
	"github.com/mholt/binding"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type newVerifyRequest struct {
	Token string
}

func (n *newVerifyRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&n.Token: binding.Field{
			Form:     "token",
			Required: true,
		},
	}
}

func (a *Account) postVerify(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
	rc := api.From(r.Context())
	nvr := &newVerifyRequest{}
	errs := binding.Bind(r, nvr)
	if errs.Len() > 0 {
		v := render.NewBadRequestView(w, rc, langs...)
		a.C.MustRender(v)
		return
	}

	m := a.C.APIQueue.Messenger()
	err := m.DoBlocking(r.Context(), func(ctx context.Context) async.CallbackFn {
		err := a.C.Users.MarkUserVerified(a.C.F.Context(r), nvr.Token)
		return func() error {
			return err
		}
	})
	if err != nil {
		a.C.MustRenderError(w, r, errors.Wrap(err, "could not verify user"), langs...)
		return
	}

	u := paths.GetPleaseVerifyURLWithSuccess(langs[0])
	http.Redirect(w, r, u.String(), http.StatusFound)
}
