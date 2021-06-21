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

	"github.com/cjslep/dharma/esi"
	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/api/esiauth"
	"github.com/cjslep/dharma/internal/api/forum"
	"github.com/cjslep/dharma/internal/api/site"
	"github.com/cjslep/dharma/internal/async"
	"github.com/cjslep/dharma/internal/config"
	"github.com/cjslep/dharma/internal/db"
	"github.com/cjslep/dharma/internal/features"
	"github.com/cjslep/dharma/internal/log"
	"github.com/cjslep/dharma/internal/render"
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
	apiQueue *async.Queue
	fedQueue *async.Queue
	features *features.Engine

	// At config-setting time
	l   *zerolog.Logger
	oac *esi.OAuth2Client
	r   *render.Renderer

	// At build routes time
	db *db.DB
	f  app.Framework

	// At start time
	startupErr error
}

func New(f *features.Engine) *FederatedApp {
	return &FederatedApp{
		apiQueue: async.NewQueue(),
		fedQueue: async.NewQueue(),
		features: f,
	}
}

func (a *FederatedApp) apiContext() *api.Context {
	return &api.Context{
		APIQueue:              a.apiQueue,
		FedQueue:              a.fedQueue,
		OAC:                   a.oac,
		L:                     a.l,
		DB:                    a.db,
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
	// TODO
	return a.startupErr
}

func (a *FederatedApp) Stop() error {
	a.fedQueue.Stop()
	a.apiQueue.Stop()
	// TODO
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
	}
}

func (a *FederatedApp) SetConfiguration(i interface{}, apc app.APCoreConfig, debug bool) error {
	c, ok := i.(*config.Config)
	if !ok {
		return errors.New("configuration is not of type *config.Config")
	}
	h := &http.Client{} // TODO
	a.oac = &esi.OAuth2Client{
		RedirectURI: "https://" + apc.Host() + esiauth.Callback,
		ClientID:    c.ClientID,
		Secret:      c.APIKey,
		Client:      h,
	}
	a.r, a.startupErr = render.New(c, debug)
	a.l = log.Logger(c.EnableConsoleLogging, c.LogDir, c.LogFile, c.NLogFiles, c.MaxMBSizeLogFiles, c.MaxDayAgeLogFiles)
	return nil
}

func (a *FederatedApp) NotFoundHandler(f app.Framework) http.Handler {
	ctx := a.apiContext()
	return http.HandlerFunc(api.MustHaveLanguageCode(
		ctx,
		func(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
			v := render.NewHTMLView(
				w,
				http.StatusNotFound,
				"status/not_found",
				nil,
				langs...)
			ctx.MustRender(v)
		}))
}

func (a *FederatedApp) MethodNotAllowedHandler(f app.Framework) http.Handler {
	ctx := a.apiContext()
	return http.HandlerFunc(api.MustHaveLanguageCode(
		ctx,
		func(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
			v := render.NewHTMLView(
				w,
				http.StatusMethodNotAllowed,
				"status/method_not_allowed",
				nil,
				langs...)
			ctx.MustRender(v)
		}))
}

func (a *FederatedApp) InternalServerErrorHandler(f app.Framework) http.Handler {
	ctx := a.apiContext()
	return http.HandlerFunc(api.MustHaveLanguageCode(
		ctx,
		func(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
			ctx.MustRenderError(w, errors.New("an internal error occured"), langs...)
		}))
}

func (a *FederatedApp) BadRequestHandler(f app.Framework) http.Handler {
	ctx := a.apiContext()
	return http.HandlerFunc(api.MustHaveLanguageCode(
		ctx,
		func(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
			v := render.NewHTMLView(
				w,
				http.StatusBadRequest,
				"status/bad_request",
				nil,
				langs...)
			ctx.MustRender(v)
		}))
}

func (a *FederatedApp) GetLoginWebHandlerFunc(f app.Framework) http.HandlerFunc {
	// TODO
	return nil
}

func (a *FederatedApp) GetAuthWebHandlerFunc(f app.Framework) http.HandlerFunc {
	// TODO
	return nil
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
		&forum.Forum{ctx},
		&site.Site{ctx},
		&esiauth.ESIAuth{ctx},
	}
	api.BuildRoutes(ar, r, ctx)
	return nil
}

func (a *FederatedApp) NewIDPath(c context.Context, t vocab.Type) (path string, err error) {
	// TODO
	return "", nil
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
	// TODO
	return nil
}

func (a *FederatedApp) DefaultUserPrivileges() interface{} {
	// TODO
	return nil
}

func (a *FederatedApp) DefaultAdminPrivileges() interface{} {
	// TODO
	return nil
}

func (a *FederatedApp) Software() app.Software {
	// TODO
	return app.Software{}
}

func (a *FederatedApp) GetInboxWebHandlerFunc(f app.Framework) func(w http.ResponseWriter, r *http.Request, outbox vocab.ActivityStreamsOrderedCollectionPage) {
	// TODO
	return nil
}

func (a *FederatedApp) ApplyFederatingCallbacks(fwc *pub.FederatingWrappedCallbacks) (others []interface{}) {
	// TODO
	return nil
}
