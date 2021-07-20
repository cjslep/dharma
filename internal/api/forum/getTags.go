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
	"github.com/go-fed/apcore/util"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

func (f *Forum) getTags(w http.ResponseWriter, r *http.Request, langs []language.Tag) {
	tag := mux.Vars(r)["tag"]
	dataTag := data.ToTag(tag)
	var tp []data.ThreadPreview
	m := f.C.APIQueue.Messenger()
	tagcb := m.DoAsync(f.C.F.Context(r), func(ctx context.Context) async.CallbackFn {
		lang := language.English
		if len(langs) > 0 {
			lang = langs[0]
		}
		l, err := f.C.Tags.GetThreadPreviewsForTag(util.Context{ctx}, dataTag, f.NListThreads, f.MaxHTMLDepth /*TODO: page=*/, 0, lang)
		return func() error {
			tp = l
			return err
		}
	})

	// TODO: Obtain avatar information for each participant

	tagdone := <-tagcb
	err := tagdone()
	if err != nil {
		f.C.MustRenderError(w, r, errors.Wrap(err, "could not obtain thread previews to render"), langs...)
		return
	}

	rc := api.From(r.Context())
	v := render.NewHTMLView(
		w,
		http.StatusOK,
		"forum/tags",
		rc,
		map[string]interface{}{
			"tag":      dataTag,
			"previews": tp,
		},
		langs...)
	f.C.MustRender(v)
}
