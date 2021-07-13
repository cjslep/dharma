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
	"github.com/cjslep/dharma/internal/render"
	"github.com/cjslep/dharma/internal/services"
	"github.com/go-fed/apcore/app"
	"github.com/go-fed/apcore/util"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

func (f *Forum) getForum(w http.ResponseWriter, r *http.Request, k app.Session, langs []language.Tag) {
	m := f.C.APIQueue.Messenger()
	var lt map[string]*services.LatestTag
	tagcb := m.DoAsync(f.C.F.Context(r), func(ctx context.Context) async.CallbackFn {
		lang := language.English
		if len(langs) > 0 {
			lang = langs[0]
		}
		l, err := f.C.Tags.GetLatestSnippets(util.Context{ctx}, f.Display, f.NPreview, f.SizePreview, f.MaxHTMLDepth, lang)
		return func() error {
			lt = l
			return err
		}
	})

	// TODO: Obtain avatar information for each participant

	tagdone := <-tagcb
	err := tagdone()
	if err != nil {
		f.C.MustRenderError(w, r, errors.Wrap(err, "could not obtain tag previews to render"), langs...)
		return
	}

	rc := api.From(r.Context())
	v := render.NewHTMLView(
		w,
		http.StatusOK,
		"forum/home",
		rc,
		map[string]interface{}{
			"preview": lt,
		},
		langs...)
	f.C.MustRender(v)
}
