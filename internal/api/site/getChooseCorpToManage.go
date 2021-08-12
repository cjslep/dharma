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
	"net/url"

	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/render"
	"golang.org/x/text/language"
)

func getChooseCorpToManageURLNotCEO(r *http.Request) *url.URL {
	return getChooseCorpToManageURL(r, "notCEO")
}

func getChooseCorpToManageURL(r *http.Request, err string) *url.URL {
	u := &url.URL{}
	u.Path = r.URL.Path
	v := url.Values{}
	v.Add("err", err)
	u.RawQuery = v.Encode()
	return u
}

func (s *Site) getChooseCorpToManage(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
	rc := api.From(r.Context())
	err := r.URL.Query().Get("err")
	v := render.NewHTMLView(
		w,
		http.StatusOK,
		"site/choose_corp_to_manage",
		rc,
		map[string]interface{}{
			"err": err,
		},
		langs...)
	s.C.MustRender(v)
}
