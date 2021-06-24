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

package render

import (
	"html/template"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/unrolled/render"
	"golang.org/x/text/language"
)

type View struct {
	w        io.Writer
	status   int
	htmlData *htmlData
	jsonData *jsonData
}

type htmlData struct {
	Name  string
	Data  map[string]interface{}
	Langs []string
}

type jsonData struct {
	Payload interface{}
}

type RenderNavDataGetter interface {
	RenderNavData() map[string]interface{}
}

func NewErrorView(w io.Writer, l *zerolog.Logger, e error, rc RenderNavDataGetter, langs ...language.Tag) *View {
	l.Error().Stack().Err(e).Msg("")
	return NewHTMLView(w, http.StatusInternalServerError, "status/internal_error", rc, map[string]interface{}{
		"err": e,
	}, langs...)
}

func NewHTMLView(w io.Writer, status int, name string, rc RenderNavDataGetter, data map[string]interface{}, langs ...language.Tag) *View {
	l := make([]string, len(langs))
	for i := range langs {
		l[i] = langs[i].String()
	}
	if data == nil {
		data = map[string]interface{}{}
	}
	data["nav"] = rc.RenderNavData()
	return &View{
		w:      w,
		status: status,
		htmlData: &htmlData{
			Name:  name,
			Data:  data,
			Langs: l,
		},
	}
}

func NewJSONView(w io.Writer, status int, payload interface{}) *View {
	return &View{
		w:      w,
		status: status,
		jsonData: &jsonData{
			Payload: payload,
		},
	}
}

func (v *View) isHTML() bool {
	return v.htmlData != nil
}

func (v *View) isJSON() bool {
	return v.jsonData != nil
}

func (v *View) render(base render.Options, funcMap func(langs ...string) []template.FuncMap) error {
	if v.isHTML() {
		base.Funcs = funcMap(v.htmlData.Langs...)
		r := render.New(base)
		return v.html(r)
	} else if v.isJSON() {
		base.Funcs = funcMap()
		r := render.New(base)
		return v.json(r)
	} else {
		return errors.New("unsupported view type")
	}
}

func (v *View) html(r *render.Render) error {
	return r.HTML(v.w, v.status, v.htmlData.Name, v.htmlData.Data)
}

func (v *View) json(r *render.Render) error {
	return r.JSON(v.w, v.status, v.jsonData.Payload)
}
