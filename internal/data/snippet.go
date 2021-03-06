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

type Snippet struct {
	ID             *url.URL
	Authors        []*url.URL
	PreviewContent string
	Created        time.Time
	Type           string
}

type snippetable interface {
	extract.IDable
	extract.Authorable
	extract.Contentable
	extract.Publishable
	GetTypeName() string
}

func ToSnippet(t vocab.Type, n, maxDepth int, preferLang language.Tag) Snippet {
	s, ok := t.(snippetable)
	if !ok {
		return Snippet{}
	}

	var p Snippet
	p.ID = extract.ToID(s)
	if p.ID == nil {
		return p
	}
	p.Authors = extract.ToAuthors(s)
	p.PreviewContent = extract.ToPreviewContent(s, n, maxDepth, preferLang)
	p.Created = extract.ToCreated(s)
	p.Type = s.GetTypeName()
	return p
}
