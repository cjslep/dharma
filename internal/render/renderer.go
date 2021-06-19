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
	"runtime"

	"github.com/cjslep/dharma/assets"
	"github.com/cjslep/dharma/internal/config"
	"github.com/cjslep/dharma/locales"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/unrolled/render"
	"golang.org/x/text/language"
)

type Renderer struct {
	b     *i18n.Bundle
	r     *render.Render
	debug bool
}

func New(c *config.Config, debug bool) (*Renderer, error) {
	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	if err := locales.AddMessageFiles(b); err != nil {
		return nil, err
	}

	r := &Renderer{
		b:     b,
		debug: debug,
	}
	rn := render.New(render.Options{
		Funcs:                     r.newFuncMap(),
		JSONContentType:           "application/json;charset=utf-8",
		Asset:                     r.asset,
		AssetNames:                r.assetNames,
		IsDevelopment:             r.debug,
		DisableHTTPErrorRendering: true,
	})
	r.r = rn

	return r, nil
}

func (r *Renderer) Render(v *View) error {
	if v.isHTML() {
		return r.html(v)
	} else if v.isJSON() {
		return r.json(v)
	} else {
		return errors.New("unsupported view type")
	}
}

func (r *Renderer) html(v *View) error {
	return r.r.HTML(v.w, v.status, v.htmlData.Name, v.htmlData.Data)
}

func (r *Renderer) json(v *View) error {
	return r.r.JSON(v.w, v.status, v.jsonData.Payload)
}

func (r *Renderer) asset(name string) ([]byte, error) {
	return assets.Asset(name)
}

func (r *Renderer) assetNames() []string {
	return assets.AssetNames()
}

func (r *Renderer) newFuncMap() []template.FuncMap {
	return []template.FuncMap{
		{
			// Application
			"goVersion": func() string {
				return runtime.Version()
			},
			"dharmaName": func() string {
				return "" // TODO
			},
			"version": func() string {
				return "" // TODO
			},
			"templateDir": func() string {
				return "" // TODO
			},
			"imgDir": func() string {
				return "" // TODO
			},
			"jsDir": func() string {
				return "" // TODO
			},
			"cssDir": func() string {
				return "" // TODO
			},
			// Debug
			"debug": func() bool {
				return r.debug
			},
			// Date & Time
			"longDate": func() string {
				return "" // TODO
			},
			"shortDate": func() string {
				return "" // TODO
			},
			// EVE Online
			"corpName": func() string {
				return "" // TODO
			},
			"allianceName": func() string {
				return "" // TODO
			},
			"coalitionAffiliations": func() []string {
				return nil // TODO
			},
		},
	}
}
