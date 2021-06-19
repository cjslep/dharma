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

package server

import (
	"net/http"

	"github.com/cjslep/dharma/esi"
	"github.com/cjslep/dharma/internal/activitypub"
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
	"github.com/go-fed/apcore/app"
	"github.com/rs/zerolog"
)

type Server struct {
	*activitypub.FederatedApp

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

func New(f *features.Engine) *Server {
	w := &Server{
		apiQueue: async.NewQueue(),
		fedQueue: async.NewQueue(),
		features: f,
	}
	w.FederatedApp = activitypub.New(w.fedQueue.Messenger(), w.start, w.stop, w.buildRoutes, w.onSetConfiguration)
	return w
}

func (w *Server) apiContext() *api.Context {
	return &api.Context{
		APIQueue:     w.apiQueue,
		FedQueue:     w.fedQueue,
		OAC:          w.oac,
		L:            w.l,
		DB:           w.db,
		F:            w.f,
		Features:     w.features,
		ErrorHandler: w.errorHandler,
	}
}

func (w *Server) buildRoutes(ar app.Router, d app.Database, f app.Framework) error {
	w.db = db.New(d)
	w.f = f
	ctx := w.apiContext()
	r := []api.Router{
		&forum.Forum{ctx},
		&site.Site{ctx},
		&esiauth.ESIAuth{ctx},
	}
	api.BuildRoutes(ar, r)
	return nil
}

func (w *Server) start() error {
	w.apiQueue.Start()
	w.fedQueue.Start()
	// TODO
	return w.startupErr
}

func (w *Server) stop() error {
	w.fedQueue.Stop()
	w.apiQueue.Stop()
	// TODO
	return nil
}

func (w *Server) onSetConfiguration(c *config.Config, apc app.APCoreConfig, debug bool) {
	h := &http.Client{} // TODO
	w.oac = &esi.OAuth2Client{
		RedirectURI: "https://" + apc.Host() + esiauth.Callback,
		ClientID:    c.ClientID,
		Secret:      c.APIKey,
		Client:      h,
	}
	w.r, w.startupErr = render.New(c, debug)
	w.l = log.Logger(c.EnableConsoleLogging, c.LogDir, c.LogFile, c.NLogFiles, c.MaxMBSizeLogFiles, c.MaxDayAgeLogFiles)
}

func (w *Server) errorHandler(wr http.ResponseWriter, r *http.Request, err error) {
	w.l.Error().Stack().Err(err).Msg("")
	w.InternalServerErrorHandler(w.f).ServeHTTP(wr, r)
}
