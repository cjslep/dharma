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
	"path"

	"github.com/cjslep/dharma/assets"
	"github.com/go-fed/apcore/app"
	"github.com/gorilla/mux"
	"golang.org/x/text/language"
)

type Router interface {
	Route(app.Router)
}

func BuildRoutes(ar app.Router, rt []Router, ctx *Context) {
	ar.Use(getPath())
	ar.Use(getSession(ctx))
	assets.AddAssetHandlers(ar)
	// Capture the locale in routing HTML rendered web pages
	ar.NewRoute().WebOnlyHandler("/", redirToEnHomepage())
	localeRouter := ar.PathPrefix("/{locale}").Subrouter()
	localeRouter.Use(getLanguageTags(ctx))
	for _, r := range rt {
		r.Route(localeRouter)
	}
}

// Sets the current request path in the request context.
func getPath() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rc := From(r.Context())
			rc.WithPath(r.URL.Path)
			r = rc.Update(r)
			next.ServeHTTP(w, r)
		})
	}
}

// Adds a session to the request context, or serves an error
// page if unable.
func getSession(ctx *Context) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			k, err := ctx.F.Session(r)
			if err != nil {
				ctx.MustRenderError(w, r, err)
				return
			}
			rc := From(r.Context())
			rc.WithSession(k)
			r = rc.Update(r)

			next.ServeHTTP(w, r)
		})
	}
}

// Sets the languages requested in the request's context:
// 1. If there is a language code in the first part of the
//    path, and it is a supported language code, then use that.
// 2. If there is a language code in the first part of the
//    path, and it is not a supported language code, redirect to
//    the closest matching supported language.
// 3. If there is no language code, use the default 'en'
//    language and redirect.
func getLanguageTags(ctx *Context) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			tag, err := language.Parse(vars["locale"])
			if err != nil {
				u := r.URL
				u.Path = path.Join("/en", u.Path)
				http.Redirect(w, r, u.String(), http.StatusFound)
				return
			}

			m := language.NewMatcher(ctx.SupportedLanguageTags())
			matchedTag, _, _ := m.Match(tag)
			if matchedTag != tag {
				u := r.URL
				u.Path = path.Join("/"+matchedTag.String(), u.Path)
				http.Redirect(w, r, u.String(), http.StatusFound)
				return
			}

			rc := From(r.Context())
			rc.WithLanguageTags(matchedTag)
			r = rc.Update(r)
			next.ServeHTTP(w, r)
		})
	}
}

func redirToEnHomepage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/en/", http.StatusFound)
	})
}
