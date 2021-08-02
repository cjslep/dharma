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
	"github.com/cjslep/dharma/internal/sessions"
	"github.com/go-fed/apcore/app"
	"github.com/go-fed/apcore/util"
	"github.com/mholt/binding"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type characterSelectRequest struct {
	CharacterID int32
}

func (c *characterSelectRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&c.CharacterID: binding.Field{
			Form:     "character_id",
			Required: true,
		},
	}
}

func (a *Account) postCharacters(w http.ResponseWriter, r *http.Request, k app.Session, langs []language.Tag) {
	rc := api.From(r.Context())
	csr := &characterSelectRequest{}
	errs := binding.Bind(r, csr)
	if errs.Len() > 0 {
		v := render.NewBadRequestView(w, rc, langs...)
		a.C.MustRender(v)
		return
	}

	userID, err := k.UserID()
	if err != nil {
		v := render.NewBadRequestView(w, rc, langs...)
		a.C.MustRender(v)
		return
	}

	if ok, err := a.C.ESI.HasCharacterForUser(util.Context{r.Context()}, userID, csr.CharacterID); err != nil {
		a.C.MustRenderError(w, r, errors.Wrap(err, "could not determine if user has character"), langs...)
		return
	} else if !ok {
		v := render.NewBadRequestView(w, rc, langs...)
		a.C.MustRender(v)
		return
	}

	sessions.SetCharacterSelected(k, csr.CharacterID)
	if err := k.Save(r, w); err != nil {
		a.C.MustRenderErrorEnglish(w, r, errors.Wrap(err, "could not save session"))
		return
	}

	// TODO: Redirect
}
