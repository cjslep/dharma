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

package activitypub

import (
	"context"
	"net/http"
	"strings"

	"github.com/cjslep/dharma/esi"
	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/api/account"
	"github.com/cjslep/dharma/internal/api/esiauth"
	"github.com/cjslep/dharma/internal/api/forum"
	"github.com/cjslep/dharma/internal/api/site"
	"github.com/cjslep/dharma/internal/async"
	"github.com/cjslep/dharma/internal/config"
	"github.com/cjslep/dharma/internal/data"
	"github.com/cjslep/dharma/internal/db"
	"github.com/cjslep/dharma/internal/features"
	"github.com/cjslep/dharma/internal/log"
	"github.com/cjslep/dharma/internal/render"
	"github.com/cjslep/dharma/internal/services"
	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/go-fed/apcore/app"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
)

var _ app.Application = new(FederatedApp)
var _ app.S2SApplication = new(FederatedApp)

type FederatedApp struct {
	// At constructor time
	software app.Software
	apiQueue *async.Queue
	fedQueue *async.Queue
	features *features.Engine

	// At config-setting time
	config *config.Config
	l      *zerolog.Logger
	oac    *esi.OAuth2Client
	r      *render.Renderer

	// At build routes time
	db *db.DB
	f  app.Framework

	// At start time
	startupErr error
}

func New(bg context.Context, f *features.Engine, software app.Software) *FederatedApp {
	return &FederatedApp{
		software: software,
		apiQueue: async.NewQueue(bg),
		fedQueue: async.NewQueue(bg),
		features: f,
	}
}

func (a *FederatedApp) apiContext() *api.Context {
	return &api.Context{
		APIQueue:              a.apiQueue,
		FedQueue:              a.fedQueue,
		OAC:                   a.oac,
		L:                     a.l,
		ESI:                   &services.ESI{a.db},
		Tags:                  &services.Tags{a.db},
		Posts:                 &services.Posts{a.db, a.f, a.fedQueue},
		Threads:               &services.Threads{a.db},
		F:                     a.f,
		Features:              a.features,
		MustRender:            a.mustRender,
		SupportedLanguageTags: a.r.LanguageTags,
	}
}

func (a *FederatedApp) mustRender(v *render.View) {
	if err := a.r.Render(v); err != nil {
		a.l.Error().Stack().Err(err).Msg("")
	}
}

func (a *FederatedApp) Start() error {
	a.apiQueue.Start()
	a.fedQueue.Start()
	return a.startupErr
}

func (a *FederatedApp) Stop() error {
	a.fedQueue.Stop()
	a.apiQueue.Stop()
	return nil
}

func (a *FederatedApp) NewConfiguration() interface{} {
	return &config.Config{
		EnableConsoleLogging: false,
		LogDir:               "./",
		LogFile:              "dharma.log",
		NLogFiles:            5,
		MaxMBSizeLogFiles:    100,
		MaxDayAgeLogFiles:    0,
		NPreview:             3,
		LenPreview:           80,
		MaxHTMLDepth:         255,
		NListThreads:         25,
	}
}

func (a *FederatedApp) SetConfiguration(i interface{}, apc app.APCoreConfig, debug bool) error {
	c, ok := i.(*config.Config)
	if !ok {
		return errors.New("configuration is not of type *config.Config")
	}
	a.config = c
	h := &http.Client{} // TODO
	a.oac = &esi.OAuth2Client{
		RedirectURI: "https://" + apc.Host() + esiauth.Callback,
		ClientID:    c.ClientID,
		Secret:      c.APIKey,
		Client:      h,
	}
	a.r, a.startupErr = render.New(c, debug, "/static")
	a.l = log.Logger(debug || c.EnableConsoleLogging, c.LogDir, c.LogFile, c.NLogFiles, c.MaxMBSizeLogFiles, c.MaxDayAgeLogFiles)
	return nil
}

func (a *FederatedApp) NotFoundHandler(f app.Framework) http.Handler {
	ctx := a.apiContext()
	return http.HandlerFunc(api.ApplyMiddleware(
		api.MustHaveLanguageCode(
			ctx,
			func(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
				rc := api.From(r.Context())
				v := render.NewHTMLView(
					w,
					http.StatusNotFound,
					"status/not_found",
					rc,
					nil,
					langs...)
				ctx.MustRender(v)
			})))
}

func (a *FederatedApp) MethodNotAllowedHandler(f app.Framework) http.Handler {
	ctx := a.apiContext()
	return http.HandlerFunc(api.ApplyMiddleware(
		api.MustHaveLanguageCode(
			ctx,
			func(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
				rc := api.From(r.Context())
				v := render.NewHTMLView(
					w,
					http.StatusMethodNotAllowed,
					"status/method_not_allowed",
					rc,
					nil,
					langs...)
				ctx.MustRender(v)
			})))
}

func (a *FederatedApp) InternalServerErrorHandler(f app.Framework) http.Handler {
	ctx := a.apiContext()
	return http.HandlerFunc(api.ApplyMiddleware(
		api.MustHaveLanguageCode(
			ctx,
			func(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
				ctx.MustRenderError(w, r, errors.New("an internal error occured"), langs...)
			})))
}

func (a *FederatedApp) BadRequestHandler(f app.Framework) http.Handler {
	ctx := a.apiContext()
	return http.HandlerFunc(api.ApplyMiddleware(
		api.MustHaveLanguageCode(
			ctx,
			func(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
				rc := api.From(r.Context())
				v := render.NewBadRequestView(w, rc, langs...)
				ctx.MustRender(v)
			})))
}

func (a *FederatedApp) GetLoginWebHandlerFunc(f app.Framework) http.HandlerFunc {
	ctx := a.apiContext()
	return http.HandlerFunc(api.ApplyMiddleware(
		api.MustHaveLanguageCode(
			ctx,
			func(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
				rc := api.From(r.Context())
				lerr := r.URL.Query().Get("login_error")
				v := render.NewHTMLView(
					w,
					http.StatusOK,
					"auth/login",
					rc,
					map[string]interface{}{
						"loginError": lerr,
					},
					langs...)
				ctx.MustRender(v)
			})))
}

func (a *FederatedApp) GetAuthWebHandlerFunc(f app.Framework) http.HandlerFunc {
	ctx := a.apiContext()
	return http.HandlerFunc(api.ApplyMiddleware(
		api.MustHaveLanguageCode(
			ctx,
			func(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
				rc := api.From(r.Context())
				v := render.NewHTMLView(
					w,
					http.StatusOK,
					"auth/auth",
					rc,
					nil,
					langs...)
				ctx.MustRender(v)
			})))
}

func (a *FederatedApp) GetOutboxWebHandlerFunc(f app.Framework) func(w http.ResponseWriter, r *http.Request, outbox vocab.ActivityStreamsOrderedCollectionPage) {
	// TODO
	return nil
}

func (a *FederatedApp) GetFollowersWebHandlerFunc(f app.Framework) (app.CollectionPageHandlerFunc, app.AuthorizeFunc) {
	// TODO
	return nil, nil
}

func (a *FederatedApp) GetFollowingWebHandlerFunc(f app.Framework) (app.CollectionPageHandlerFunc, app.AuthorizeFunc) {
	// TODO
	return nil, nil
}

func (a *FederatedApp) GetLikedWebHandlerFunc(f app.Framework) (app.CollectionPageHandlerFunc, app.AuthorizeFunc) {
	// TODO
	return nil, nil
}

func (a *FederatedApp) GetUserWebHandlerFunc(f app.Framework) (app.VocabHandlerFunc, app.AuthorizeFunc) {
	// TODO
	return nil, nil
}

func (a *FederatedApp) BuildRoutes(ar app.Router, d app.Database, f app.Framework) error {
	a.db = db.New(d)
	a.f = f
	ctx := a.apiContext()
	r := []api.Router{
		&forum.Forum{ctx, data.AllTags, a.config.NPreview, a.config.LenPreview, a.config.MaxHTMLDepth, a.config.NListThreads},
		&site.Site{ctx},
		&account.Account{ctx},
		&esiauth.ESIAuth{ctx},
	}
	api.BuildRoutes(ar, r, ctx)
	return nil
}

func (a *FederatedApp) Paths() app.Paths {
	homepageFn := func(cur string) string {
		s := strings.Split(cur, "/")
		if len(s) < 2 {
			return "/en"
		}
		return "/" + s[1]
	}
	return app.Paths{
		GetLogin:            "/{locale}/login",
		PostLogin:           "/{locale}/login",
		GetLogout:           "/{locale}/logout",
		GetOAuth2Authorize:  "/{locale}/oauth2/authorize",
		PostOAuth2Authorize: "/{locale}/oauth2/authorize",
		RedirectToHomepage:  homepageFn,
		RedirectToLogin: func(cur string) string {
			return homepageFn(cur) + "/login"
		},
	}
}

func (a *FederatedApp) StaticServingEnabled() bool {
	// We compile assets into our binary
	return false
}

func (a *FederatedApp) NewIDPath(c context.Context, t vocab.Type) (path string, err error) {
	// TODO
	return "", errors.Errorf("unhandled type name: %s", t.GetTypeName())
}

func (a *FederatedApp) ScopePermitsPrivateGetInbox(scope string) (permitted bool, err error) {
	// TODO
	return false, nil
}

func (a *FederatedApp) ScopePermitsPrivateGetOutbox(scope string) (permitted bool, err error) {
	// TODO
	return false, nil
}

func (a *FederatedApp) DefaultUserPreferences() interface{} {
	return data.DefaultPreferences()
}

func (a *FederatedApp) DefaultUserPrivileges() interface{} {
	return data.DefaultPrivileges()
}

func (a *FederatedApp) DefaultAdminPrivileges() interface{} {
	return data.DefaultAdminPrivileges()
}

func (a *FederatedApp) Software() app.Software {
	return a.software
}

func (a *FederatedApp) GetInboxWebHandlerFunc(f app.Framework) func(w http.ResponseWriter, r *http.Request, outbox vocab.ActivityStreamsOrderedCollectionPage) {
	// TODO
	return nil
}

func (a *FederatedApp) ApplyFederatingCallbacks(fwc *pub.FederatingWrappedCallbacks) (others []interface{}) {
	// TODO
	return nil
}
