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

package forum

import (
	"net/http"

	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/render"
	"github.com/go-fed/apcore/app"
	"github.com/mholt/binding"
	"golang.org/x/text/language"
)

type newMarkdownPreviewRequest struct {
	Body string
}

func (n *newMarkdownPreviewRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&n.Body: binding.Field{
			Form:     "body",
			Required: true,
		},
	}
}

type newMarkdownPreviewResponse struct {
	Markdown string `json:"markdown"`
}

func (f *Forum) postPreviewMarkdown(w http.ResponseWriter, r *http.Request, k app.Session, langs []language.Tag) {
	rc := api.From(r.Context())
	npr := &newMarkdownPreviewRequest{}
	errs := binding.Bind(r, npr)
	if errs.Len() > 0 {
		v := render.NewBadRequestView(w, rc, langs...)
		f.C.MustRender(v)
		return
	}
	resp := &newMarkdownPreviewResponse{
		Markdown: "", // TODO: Process markdown
	}
	v := render.NewJSONView(w, http.StatusOK, resp)
	f.C.MustRender(v)
}
