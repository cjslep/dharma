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
	"net/url"

	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/async"
	"github.com/cjslep/dharma/internal/data"
	"github.com/cjslep/dharma/internal/render"
	"github.com/go-fed/apcore/app"
	"github.com/mholt/binding"
	"golang.org/x/text/language"
)

type newPostRequest struct {
	Title string
	Body  string
	Tags  []string
}

func (n *newPostRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&n.Title: binding.Field{
			Form:     "title",
			Required: true,
		},
		&n.Body: binding.Field{
			Form:     "body",
			Required: true,
		},
		&n.Tags: binding.Field{
			Form:     "tags",
			Required: true,
		},
	}
}

func (f *Forum) postNewPost(w http.ResponseWriter, r *http.Request, k app.Session, langs []language.Tag) {
	rc := api.From(r.Context())
	npr := &newPostRequest{}
	errs := binding.Bind(r, npr)
	if errs.Len() > 0 {
		v := render.NewBadRequestView(w, rc, langs...)
		f.C.MustRender(v)
		return
	}

	// TODO: Check tags for bad request
	var validatedTags []data.Tag

	// TODO: Determine the language
	var lang language.Tag

	user, err := k.UserID()
	if err != nil {
		v := render.NewBadRequestView(w, rc, langs...)
		f.C.MustRender(v)
		return
	}

	m := f.C.APIQueue.Messenger()
	var id *url.URL
	err = m.DoBlocking(r.Context(), func(ctx context.Context) async.CallbackFn {
		var err error
		id, err = f.C.Posts.CreateNewPost(ctx, npr.Title, npr.Body, user, validatedTags, lang)
		return func() error {
			return err
		}
	})

	// TODO: Redirect, bad request, internal error, etc
}
