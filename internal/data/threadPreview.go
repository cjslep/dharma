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

package data

import (
	"net/url"
	"time"

	"github.com/cjslep/dharma/internal/data/extract"
	"github.com/go-fed/activity/streams/vocab"
	"golang.org/x/text/language"
)

type ThreadPreviewMessage struct {
	ID             *url.URL
	Authors        []*url.URL
	PreviewContent string
	Created        time.Time
	Type           string
}

type ThreadPreview struct {
	Title string
	First *ThreadPreviewMessage
	Last  *ThreadPreviewMessage
}

type threadPreviewable interface {
	extract.IDable
	extract.Authorable
	extract.Contentable
	extract.Publishable
	GetTypeName() string
}

func ToThreadPreview(first, last vocab.Type, n, maxDepth int, lang language.Tag) ThreadPreview {
	tp := ThreadPreview{}
	// TODO: Title
	ftp, ok := first.(threadPreviewable)
	if ok {
		tp.First = &ThreadPreviewMessage{}
		toThreadPreviewMessage(tp.First, ftp, n, maxDepth, lang)
	}
	ltp, ok := last.(threadPreviewable)
	if ok {
		tp.Last = &ThreadPreviewMessage{}
		toThreadPreviewMessage(tp.Last, ltp, n, maxDepth, lang)
	}
	return tp
}

func toThreadPreviewMessage(tpm *ThreadPreviewMessage, t threadPreviewable, n, maxDepth int, lang language.Tag) {
	tpm.ID = extract.ToID(t)
	if tpm.ID == nil {
		return
	}
	tpm.Authors = extract.ToAuthors(t)
	tpm.PreviewContent = extract.ToPreviewContent(t, n, maxDepth, lang)
	tpm.Created = extract.ToCreated(t)
	tpm.Type = t.GetTypeName()
}
