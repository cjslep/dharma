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

package forum

import (
	"context"
	"net/http"

	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/async"
	"github.com/cjslep/dharma/internal/data"
	"github.com/cjslep/dharma/internal/render"
	"github.com/go-fed/apcore/app"
	"github.com/go-fed/apcore/util"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

func (f *Forum) getThreads(w http.ResponseWriter, r *http.Request, k app.Session, langs []language.Tag) {
	tid := mux.Vars(r)["thread"]
	var ps []data.Post
	m := f.C.APIQueue.Messenger()
	threadcb := m.DoAsync(f.C.F.Context(r), func(ctx context.Context) async.CallbackFn {
		lang := language.English
		if len(langs) > 0 {
			lang = langs[0]
		}
		t, err := f.C.Threads.GetPosts(util.Context{ctx}, tid /*TODO: n=*/, 25 /*TODO: page=*/, 0, lang)
		return func() error {
			ps = t
			return err
		}
	})

	// TODO: Obtain avatar information for each participant, here or in services

	done := <-threadcb
	err := done()
	if err != nil {
		f.C.MustRenderError(w, r, errors.Wrap(err, "could not obtain thread to render"), langs...)
		return
	}

	rc := api.From(r.Context())
	v := render.NewHTMLView(
		w,
		http.StatusOK,
		"forum/threads",
		rc,
		map[string]interface{}{
			"posts": ps,
		},
		langs...)
	f.C.MustRender(v)
}
