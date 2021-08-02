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

package site

import (
	"net/http"

	"github.com/cjslep/dharma/internal/api"
	"golang.org/x/text/language"
	"github.com/mholt/binding"
	"github.com/cjslep/dharma/internal/render"
	"github.com/go-fed/apcore/util"
	"github.com/go-fed/apcore/app"
)

type chooseCorpRequest struct {
	CorporationID int32
}

func (c *chooseCorpRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&c.CorporationID: binding.Field{
			Form:     "corporation_id",
			Required: true,
		},
	}
}

func (s *Site) postChooseCorpToManage(w http.ResponseWriter, r *http.Request, k app.Session, langs []language.Tag) {
	rc := api.From(r.Context())
	ccr := &chooseCorpRequest{}
	errs := binding.Bind(r, ccr)
	if errs.Len() > 0 {
		v := render.NewBadRequestView(w, rc, langs...)
		s.C.MustRender(v)
		return
	}

	userID, err := k.UserID()
	if err != nil {
		v := render.NewBadRequestView(w, rc, langs...)
		s.C.MustRender(v)
		return
	}

	err = s.C.State.ChooseCorporation(util.Context{r.Context()}, userID, ccr.CorporationID)
	if err != nil {
		s.C.MustRenderError(w, r, errors.Wrap(err, "could not choose corporation"), langs...)
		return
	}

	// TODO: Redirect to success message?
}
