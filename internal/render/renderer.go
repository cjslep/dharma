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
	"fmt"
	"html/template"
	"runtime"

	"github.com/cjslep/dharma/assets"
	"github.com/cjslep/dharma/internal/config"
	d_i18n "github.com/cjslep/dharma/internal/render/i18n"
	"github.com/cjslep/dharma/locales"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"github.com/unrolled/render"
	"golang.org/x/text/language"
)

type Renderer struct {
	b          *i18n.Bundle
	baseOpts   render.Options
	debug      bool
	staticRoot string
}

func New(c *config.Config, debug bool, staticRoot string) (*Renderer, error) {
	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	if err := locales.AddMessageFiles(b); err != nil {
		return nil, err
	}

	r := &Renderer{
		b:          b,
		debug:      debug,
		staticRoot: staticRoot,
	}
	r.baseOpts = render.Options{
		JSONContentType:           "application/json;charset=utf-8",
		Asset:                     r.asset,
		AssetNames:                r.assetNames,
		IsDevelopment:             r.debug,
		DisableHTTPErrorRendering: true,
	}

	return r, nil
}

func (r *Renderer) LanguageTags() []language.Tag {
	return r.b.LanguageTags()
}

func (r *Renderer) Render(v *View) error {
	return v.render(r.baseOpts, r.newFuncMap)
}

func (r *Renderer) asset(name string) ([]byte, error) {
	return assets.Asset(name)
}

func (r *Renderer) assetNames() []string {
	return assets.AssetNames()
}

func (r *Renderer) newFuncMap(langs ...string) []template.FuncMap {
	m := d_i18n.New(r.b, langs...)
	return []template.FuncMap{
		{
			// Application
			"GoVersion": func() string {
				return runtime.Version()
			},
			"DharmaName": func() string {
				return "" // TODO
			},
			"Version": func() string {
				return "" // TODO
			},
			"TemplateDir": func() string {
				return "" // TODO
			},
			"ImgDir": func() string {
				return "" // TODO
			},
			"JsDir": func() string {
				return fmt.Sprintf("%s/js", r.staticRoot)
			},
			"CssDir": func() string {
				return fmt.Sprintf("%s/css", r.staticRoot)
			},
			// Utility
			"Escape": func(s string) template.HTML {
				return template.HTML(s)
			},
			"Locale": func() *d_i18n.Messages {
				return m
			},
			// Debug
			"Debug": func() bool {
				return r.debug
			},
			// Date & Time
			"LongDate": func() string {
				return "" // TODO
			},
			"ShortDate": func() string {
				return "" // TODO
			},
			// EVE Online
			"CorpName": func() string {
				return "" // TODO: I actually need this.
			},
			"AllianceName": func() string {
				return "" // TODO
			},
			"CoalitionAffiliations": func() []string {
				return nil // TODO
			},
		},
	}
}
