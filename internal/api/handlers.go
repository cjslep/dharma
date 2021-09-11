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

func ApplyMiddleware(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getPath()(
			getSession(ctx)(
				getPrivileges(ctx)(
					getLanguageTags(ctx)(
						enforceEmailValidation(ctx)(next))))).ServeHTTP(w, r)
	})
}

func MustHaveSession(ctx *Context, r StatefulRenderHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rc := From(req.Context())
		k, err := rc.Session()
		if err != nil {
			ctx.MustRenderErrorEnglish(w, req, errors.New("could not retrieve session"))
			return
		}

		r(w, req, k)
	})
}

func MustHaveLanguageCode(r LocalizedRenderHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rc := From(req.Context())
		langs, err := rc.LanguageTags()
		if err != nil {
			langs = []language.Tag{language.English}
		}
		r(w, req, langs)
	})
}

func MustHaveSessionAndLanguageCode(ctx *Context, r LocalizedStatefulRenderHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rc := From(req.Context())
		langs, err := rc.LanguageTags()
		if err != nil {
			ctx.MustRenderErrorEnglish(w, req, errors.New("could not retrieve request language tags"))
			return
		}

		k, err := rc.Session()
		if err != nil {
			ctx.MustRenderError(w, req, errors.New("could not retrieve session"), langs...)
			return
		}

		r(w, req, k, langs)
	})
}

func CorpMustBeManaged(ctx *Context, next http.Handler) http.Handler {
	return enforceCorpIsManaged(ctx)(next)
}

func CorpMustNotBeManaged(ctx *Context, next http.Handler) http.Handler {
	return enforceCorpIsNotManaged(ctx)(next)
}

func MustBeAdmin(ctx *Context, next http.Handler) http.Handler {
	return enforceLoggedInAsAdmin(ctx)(next)
}

// TODO: Use this function
func MustHaveCharacterSelected(ctx *Context, next http.Handler) http.Handler {
	return enforceCharacterSelected(ctx)(next)
}
