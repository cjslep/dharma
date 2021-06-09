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
	"github.com/cjslep/dharma/internal/activitypub"
	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/api/forum"
	"github.com/cjslep/dharma/internal/api/site"
	"github.com/cjslep/dharma/internal/async"
	"github.com/cjslep/dharma/internal/db"
	"github.com/go-fed/apcore/app"
)

type Server struct {
	*activitypub.FederatedApp

	apiQueue *async.Queue
	fedQueue *async.Queue
	db       *db.DB
	f        app.Framework
}

func New() *Server {
	// TODO
	w := &Server{
		apiQueue: async.NewQueue(),
		fedQueue: async.NewQueue(),
	}
	w.FederatedApp = activitypub.New(w.fedQueue.Messenger(), w.start, w.stop, w.buildRoutes)
	return w
}

func (w *Server) buildRoutes(ar app.Router, d app.Database, f app.Framework) error {
	w.db = db.New(d)
	r := []api.Router{
		forum.New(w.db, w.apiQueue.Messenger()),
		site.New(w.db, w.apiQueue.Messenger()),
	}
	api.BuildRoutes(ar, r)
	return nil
}

func (w *Server) start() error {
	w.apiQueue.Start()
	w.fedQueue.Start()
	// TODO
	return nil
}

func (w *Server) stop() error {
	w.fedQueue.Stop()
	w.apiQueue.Stop()
	// TODO
	return nil
}
