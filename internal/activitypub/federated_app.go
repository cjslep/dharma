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

	"github.com/cjslep/dharma/internal/async"
	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/go-fed/apcore/app"
)

var _ app.Application = new(FederatedApp)
var _ app.S2SApplication = new(FederatedApp)

type FederatedApp struct {
	m           *async.Messenger
	start       func() error
	stop        func() error
	buildRoutes func(r app.Router, db app.Database, f app.Framework) error
}

func New(m *async.Messenger, start, stop func() error, buildRoutes func(r app.Router, db app.Database, f app.Framework) error) *FederatedApp {
	return &FederatedApp{
		m:           m,
		start:       start,
		stop:        stop,
		buildRoutes: buildRoutes,
	}
}

func (a *FederatedApp) Start() error {
	return a.start()
}

func (a *FederatedApp) Stop() error {
	return a.stop()
}

func (a *FederatedApp) NewConfiguration() interface{} {
	// TODO
	return nil
}

func (a *FederatedApp) SetConfiguration(i interface{}) error {
	// TODO
	return nil
}

func (a *FederatedApp) NotFoundHandler(f app.Framework) http.Handler {
	// TODO
	return nil
}

func (a *FederatedApp) MethodNotAllowedHandler(f app.Framework) http.Handler {
	// TODO
	return nil
}

func (a *FederatedApp) InternalServerErrorHandler(f app.Framework) http.Handler {
	// TODO
	return nil
}

func (a *FederatedApp) BadRequestHandler(f app.Framework) http.Handler {
	// TODO
	return nil
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

func (a *FederatedApp) BuildRoutes(r app.Router, db app.Database, f app.Framework) error {
	return a.buildRoutes(r, db, f)
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
