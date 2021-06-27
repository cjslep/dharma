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

package api

import (
	"io"
	"net/http"

	"github.com/cjslep/dharma/esi"
	"github.com/cjslep/dharma/internal/async"
	"github.com/cjslep/dharma/internal/features"
	"github.com/cjslep/dharma/internal/render"
	"github.com/cjslep/dharma/internal/services"
	"github.com/go-fed/apcore/app"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
)

type Context struct {
	APIQueue              *async.Queue
	FedQueue              *async.Queue
	OAC                   *esi.OAuth2Client
	L                     *zerolog.Logger
	ESI                   *services.ESI
	Tags                  *services.Tags
	Posts                 *services.Posts
	F                     app.Framework
	Features              *features.Engine
	MustRender            func(*render.View)
	SupportedLanguageTags func() []language.Tag
}

func (c *Context) MustRenderErrorEnglish(w io.Writer, r *http.Request, err error) {
	rc := From(r.Context())
	c.MustRender(render.NewErrorView(w, c.L, err, rc, language.English))
}

func (c *Context) MustRenderError(w io.Writer, r *http.Request, err error, langs ...language.Tag) {
	rc := From(r.Context())
	c.MustRender(render.NewErrorView(w, c.L, err, rc, langs...))
}
