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

type Post struct {
	ID        *url.URL
	Authors   []*url.URL
	Content   string
	Created   time.Time
	Type      string
	InReplyTo []*url.URL
}

type postable interface {
	extract.IDable
	extract.Authorable
	extract.Contentable
	extract.Publishable
	extract.InReplyToable
	GetTypeName() string
}

func ToPost(t vocab.Type, lang language.Tag) Post {
	p, ok := t.(postable)
	if !ok {
		return Post{}
	}
	var o Post
	o.ID = extract.ToID(p)
	if o.ID == nil {
		return o
	}
	o.Authors = extract.ToAuthors(p)
	o.Content = extract.ToContent(p, lang)
	o.Created = extract.ToCreated(p)
	o.InReplyTo = extract.ToInReplyTo(p)
	o.Type = p.GetTypeName()
	return o
}
