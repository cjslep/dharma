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
	"net/http"

	"github.com/go-fed/apcore/app"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type StatefulRenderHandler func(w http.ResponseWriter, r *http.Request, k app.Session)
type LocalizedRenderHandler func(w http.ResponseWriter, r *http.Request, langs []language.Tag)
type LocalizedStatefulRenderHandler func(w http.ResponseWriter, r *http.Request, k app.Session, langs []language.Tag)

func MustHaveSession(ctx *Context, r StatefulRenderHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		rc := From(req.Context())
		k, err := rc.Session()
		if err != nil {
			ctx.MustRenderErrorEnglish(w, errors.New("could not retrieve session"))
			return
		}

		r(w, req, k)
	}
}

func MustHaveLanguageCode(ctx *Context, r LocalizedRenderHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		rc := From(req.Context())
		langs, err := rc.LanguageTags()
		if err != nil {
			ctx.MustRenderErrorEnglish(w, errors.New("could not retrieve request language tags"))
			return
		}
		r(w, req, langs)
	}
}

func MustHaveSessionAndLanguageCode(ctx *Context, r LocalizedStatefulRenderHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		rc := From(req.Context())
		langs, err := rc.LanguageTags()
		if err != nil {
			ctx.MustRenderErrorEnglish(w, errors.New("could not retrieve request language tags"))
			return
		}

		k, err := rc.Session()
		if err != nil {
			ctx.MustRenderError(w, errors.New("could not retrieve session"), langs...)
			return
		}

		r(w, req, k, langs)
	}
}
