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

package services

import (
	"context"
	"net/url"
	"time"

	"github.com/cjslep/dharma/internal/async"
	"github.com/cjslep/dharma/internal/data"
	"github.com/cjslep/dharma/internal/db"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/apcore/app"
	"github.com/go-fed/apcore/paths"
	"github.com/go-fed/apcore/util"
	"golang.org/x/text/language"
)

type Posts struct {
	DB    *db.DB
	F     app.Framework
	Queue *async.Queue
}

func (p *Posts) CreateNewPost(c context.Context, title, body, user string, tags []data.Tag, lang language.Tag) (*url.URL, error) {
	// TODO: Verify body is well-formed markdown
	// TODO: Determine appropraite audience to send this note to.

	// Produce the ActivityStreams data for this interaction
	// note ActivityStreams
	note := streams.NewActivityStreamsNote()

	// 'to' property
	toP := streams.NewActivityStreamsToProperty()
	// TODO
	note.SetActivityStreamsTo(toP)

	// 'summary' property TODO: Should it be 'name' instead, for a 'title'?
	summaryP := streams.NewActivityStreamsSummaryProperty()
	summaryP.AppendXMLSchemaString(title)
	// TODO: 'summary' mapping property
	note.SetActivityStreamsSummary(summaryP)

	// 'content' property
	contentP := streams.NewActivityStreamsContentProperty()
	contentP.AppendXMLSchemaString(body)
	// TODO: 'content' mapping property
	note.SetActivityStreamsContent(contentP)

	// 'mediaType' property
	mediaTypeP := streams.NewActivityStreamsMediaTypeProperty()
	mediaTypeP.Set("text/markdown")
	note.SetActivityStreamsMediaType(mediaTypeP)

	// 'published' property
	publishedP := streams.NewActivityStreamsPublishedProperty()
	publishedP.Set(time.Now())
	note.SetActivityStreamsPublished(publishedP)

	// 'tag' property
	tagP := streams.NewActivityStreamsTagProperty()
	// TODO
	note.SetActivityStreamsTag(tagP)

	// Federate the ActivityStreams data
	m := p.Queue.Messenger()
	var noteIRI *url.URL
	err := m.DoBlocking(c, func(ctx context.Context) async.CallbackFn {
		err := p.F.Send(util.Context{ctx}, paths.UUID(user), note)
		return func() error {
			if err == nil {
				noteIRI = note.GetJSONLDId().GetIRI()
			}
			return err
		}
	})
	return noteIRI, err
}
