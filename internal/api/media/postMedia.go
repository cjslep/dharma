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

package media

import (
	"bytes"
	"mime/multipart"
	"net/http"

	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/render"
	"github.com/go-fed/apcore/app"
	"github.com/mholt/binding"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type newMediaRequest struct {
	Media *multipart.FileHeader
}

func (n *newMediaRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&n.Media: binding.Field{
			Form:     "media",
			Required: true,
		},
	}
}

// TODO: Include privilege information when adding media
func (m *Media) postMedia(w http.ResponseWriter, r *http.Request, k app.Session, langs []language.Tag) {
	rc := api.From(r.Context())
	nmr := &newMediaRequest{}
	errs := binding.Bind(r, nmr)
	if errs.Len() > 0 {
		v := render.NewBadRequestView(w, rc, langs...)
		m.C.MustRender(v)
		return
	}

	fh, err := nmr.Media.Open()
	if err != nil {
		m.C.MustRenderError(w, r, errors.Wrap(err, "could not open multipart form"), langs...)
		return
	}
	defer fh.Close()

	var data bytes.Buffer
	size, err := data.ReadFrom(fh)
	if err != nil {
		m.C.MustRenderError(w, r, errors.Wrap(err, "could not read from multipart form"), langs...)
		return
	} else if size > m.MaxSize {
		v := render.NewBadRequestView(w, rc, langs...)
		m.C.MustRender(v)
		return
	}
	b := data.Bytes()

	id, err := m.C.Media.CreateMedia(r.Context(), nmr.Media.Filename, http.DetectContentType(b), b)
	if err != nil {
		m.C.MustRenderError(w, r, errors.Wrap(err, "could not create media"), langs...)
		return
	}

	v := render.NewJSONView(w, http.StatusOK, struct {
		ID string `json:"id"`
	}{
		ID: id,
	})
	m.C.MustRender(v)
}
