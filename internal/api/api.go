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
	"github.com/cjslep/dharma/internal/api/paths"
	"github.com/cjslep/dharma/internal/render"
	"github.com/cjslep/dharma/internal/services"
	"github.com/cjslep/dharma/internal/sessions"
	"github.com/go-fed/apcore/app"
	ap_paths "github.com/go-fed/apcore/paths"
	"github.com/go-fed/apcore/util"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type Router interface {
	Route(app.Router)
}

func BuildRoutes(ar app.Router, rt []Router, ctx *Context) {
	ar.Use(getPath())
	ar.Use(getSession(ctx))
	ar.Use(getPrivileges(ctx))
	assets.AddAssetHandlers(ar)
	// Capture the locale in routing HTML rendered web pages
	ar.NewRoute().WebOnlyHandler("/", redirToEnHomepage())
	localeRouter := ar.PathPrefix("/{locale}").Subrouter()
	localeRouter.Use(getLanguageTags(ctx))
	localeRouter.Use(enforceEmailValidation(ctx))
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

// Adds the request's privileges and admin status.
//
// Requires a session to function.
func getPrivileges(ctx *Context) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rc := From(r.Context())
			k, err := rc.Session()
			if err == nil {
				var priv services.Privileges
				userID, err := k.UserID()
				if err == nil {
					admin, err := ctx.F.GetPrivileges(util.Context{rc}, ap_paths.UUID(userID), &priv)
					if err != nil {
						admin = false
						priv = services.DefaultPrivileges()
					}
					rc.WithIsAdmin(admin)
					rc.WithPrivileges(priv)
					r = rc.Update(r)
				}
			}

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

// If a user is signed in but has not verified their email, redirect
// to a placeholder page.
func enforceEmailValidation(ctx *Context) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rc := From(r.Context())
			k, err := rc.Session()
			if err != nil {
				// No session: abort check
				next.ServeHTTP(w, r)
				return
			}
			userID, err := k.UserID()
			if err != nil {
				// No user ID: success
				next.ServeHTTP(w, r)
				return
			}
			valid, err := ctx.Users.IsUserVerified(ctx.F.Context(r), userID)
			if err != nil {
				// Error
				ctx.MustRenderError(w, r, err)
				return
			} else if !valid && !paths.IsVerifyPath(r.URL.Path) {
				// Not yet validated & not already the
				// verify path: redirect
				lts, err := rc.LanguageTags()
				if err != nil {
					ctx.MustRenderError(w, r, err)
					return
				}
				u := paths.GetPleaseVerifyURLWithTY(lts[0])
				http.Redirect(w, r, u.String(), http.StatusFound)
				return
			}
			// Valid & no error, or the verify page
			next.ServeHTTP(w, r)
		})
	}
}

// enforceCorpIsManaged ensures that the current state of the application is
// managing a corporation before allowing HTTP requests to proceed.
func enforceCorpIsManaged(ctx *Context) mux.MiddlewareFunc {
	return enforceCorpManagedState(ctx, true)
}

// enforceCorpIsManaged ensures that the current state of the application is
// not managing a corporation before allowing HTTP requests to proceed.
func enforceCorpIsNotManaged(ctx *Context) mux.MiddlewareFunc {
	return enforceCorpManagedState(ctx, false)
}

// enforceCorpManagedState is a helper for enforcing either "must be managed" or
// "must not be managed"
func enforceCorpManagedState(ctx *Context, mustBeState bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if mustBeState && ctx.State.RequiresCorpToBeManaged() {
				rc := From(r.Context())
				langs, err := rc.LanguageTags()
				if err != nil {
					langs = []language.Tag{language.English}
				}
				v := render.NewHTMLView(
					w,
					http.StatusOK,
					"state/manage_corp",
					rc,
					map[string]interface{}{},
					langs...)
				ctx.MustRender(v)
				return
			} else if !mustBeState && !ctx.State.RequiresCorpToBeManaged() {
			}
			next.ServeHTTP(w, r)
		})
	}
}

// enforceLoggedInAsAdmin ensures that the request is for a logged-in
// user with administrative privileges.
func enforceLoggedInAsAdmin(ctx *Context) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rc := From(r.Context())
			admin, err := rc.IsAdmin()
			if err != nil {
				ctx.MustRenderError(w, r, err)
				return
			}
			if !admin {
				langs, err := rc.LanguageTags()
				if err != nil {
					langs = []language.Tag{language.English}
				}
				ctx.MustRender(render.NewNotFoundView(w, rc, langs...))
			}
			next.ServeHTTP(w, r)
		})
	}
}

// enforceCharacterSelected ensures that the endpoint is hit only if the user
// has selected an active character to use, and that character's token is not
// flagged as needing a re-scoping.
func enforceCharacterSelected(ctx *Context) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rc := From(r.Context())
			k, err := rc.Session()
			if err != nil {
				ctx.MustRenderError(w, r, errors.New("could not obtain session for enforcing character selection"))
				return
			}

			charID := sessions.GetCharacterSelected(k)
			if charID == 0 {
				// No character selected
				lts, err := rc.LanguageTags()
				if err != nil {
					ctx.MustRenderError(w, r, err)
					return
				}
				u := paths.GetCharacterSelection(lts[0])
				http.Redirect(w, r, u.String(), http.StatusFound)
				return
			}

			rescope, err := ctx.ESI.DoesCharacterNeedRescope(ctx.F.Context(r), charID)
			if err != nil {
				ctx.MustRenderError(w, r, err)
				return
			} else if rescope {
				lts, err := rc.LanguageTags()
				if err != nil {
					ctx.MustRenderError(w, r, err)
					return
				}
				u := paths.GetESIAuthPathRescope(lts[0])
				http.Redirect(w, r, u.String(), http.StatusFound)
				return
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
