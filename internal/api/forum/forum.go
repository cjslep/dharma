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
	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/data"
	"github.com/go-fed/apcore/app"
)

type Forum struct {
	C            *api.Context
	Display      []data.Tag
	NPreview     int
	SizePreview  int
	MaxHTMLDepth int
	NListThreads int
}

func (f *Forum) Route(r app.Router) {
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/forum",
		api.CorpMustBeManaged(f.C,
			api.MustHaveLanguageCode(f.C, f.getForum)))
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/forum/tags/{tag}",
		api.CorpMustBeManaged(f.C,
			api.MustHaveLanguageCode(f.C, f.getTags)))
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/forum/threads/{thread}",
		api.CorpMustBeManaged(f.C,
			api.MustHaveLanguageCode(f.C, f.getThreads)))
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/forum/posts/new",
		api.CorpMustBeManaged(f.C,
			api.MustHaveLanguageCode(f.C, f.getNewPost)))
	r.NewRoute().Methods("POST").WebOnlyHandler(
		"/forum/posts/new",
		api.CorpMustBeManaged(f.C,
			api.MustHaveSessionAndLanguageCode(f.C, f.postNewPost)))
	r.NewRoute().Methods("POST").WebOnlyHandler(
		"/forum/posts/preview/markdown",
		api.CorpMustBeManaged(f.C,
			api.MustHaveLanguageCode(f.C, f.postPreviewMarkdown)))
}
