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
	"strings"

	"github.com/cjslep/dharma/esi"
	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/render"
	"github.com/mholt/binding"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type newCorpSetupSearch struct {
	Query string
}

func (n *newCorpSetupSearch) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&n.Query: binding.Field{
			Form:     "query",
			Required: true,
		},
	}
}

func (s *Site) postCorpSetupSearch(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
	rc := api.From(r.Context())
	nse := &newCorpSetupSearch{}
	errs := binding.Bind(r, nse)
	if errs.Len() > 0 {
		v := render.NewBadRequestView(w, rc, langs...)
		s.C.MustRender(v)
		return
	}

	var lang language.Tag
	if len(langs) > 0 {
		lang = langs[0]
	} else {
		v := render.NewBadRequestView(w, rc, langs...)
		s.C.MustRender(v)
		return
	}

	nse.Query = strings.TrimSpace(nse.Query)
	if len(nse.Query) <= 3 {
		v := render.NewBadRequestView(w, rc, langs...)
		s.C.MustRender(v)
		return
	}

	corps, err := s.C.ESI.SearchCorporations(r.Context(), nse.Query, lang)
	if err != nil {
		s.C.MustRenderError(w, r, errors.Wrap(err, "could not search corporations"), langs...)
		return
	}

	resp := struct {
		Corporations []*esi.Corporation `json:"corporations"`
	}{
		Corporations: corps,
	}

	v := render.NewJSONView(w, http.StatusOK, resp)
	s.C.MustRender(v)
}
