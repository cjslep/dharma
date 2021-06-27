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

	"github.com/cjslep/dharma/internal/util"
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
	GetJSONLDId() vocab.JSONLDIdProperty
	GetActivityStreamsAttributedTo() vocab.ActivityStreamsAttributedToProperty
	GetActivityStreamsContent() vocab.ActivityStreamsContentProperty
	GetActivityStreamsPublished() vocab.ActivityStreamsPublishedProperty
	GetTypeName() string
}

func ToSnippet(t vocab.Type, n, maxDepth int, preferLang language.Tag) Snippet {
	s, ok := t.(snippetable)
	if !ok {
		return Snippet{}
	}

	var p Snippet
	idp := s.GetJSONLDId()
	if idp == nil {
		return p
	}
	p.ID = idp.GetIRI()

	atp := s.GetActivityStreamsAttributedTo()
	if atp != nil {
		for iter := atp.Begin(); iter != atp.End(); iter = iter.Next() {
			if iter.IsIRI() {
				p.Authors = append(p.Authors, iter.GetIRI())
			} else if iter.IsActivityStreamsPerson() {
				pidp := iter.GetActivityStreamsPerson().GetJSONLDId()
				if pidp != nil {
					p.Authors = append(p.Authors, pidp.GetIRI())
				}
			}
		}
	}

	cp := s.GetActivityStreamsContent()
	if cp != nil {
		for iter := cp.Begin(); iter != cp.End(); iter = iter.Next() {
			if iter.IsXMLSchemaString() {
				p.PreviewContent = util.FirstNHTML(iter.GetXMLSchemaString(), n, maxDepth)
			} else if iter.IsRDFLangString() {
				if iter.HasLanguage(preferLang.String()) {
					p.PreviewContent = util.FirstNHTML(iter.GetLanguage(preferLang.String()), n, maxDepth)
				} else if iter.HasLanguage("en") {
					p.PreviewContent = util.FirstNHTML(iter.GetLanguage("en"), n, maxDepth)
				}
			}
		}
	}

	pp := s.GetActivityStreamsPublished()
	if pp != nil {
		if pp.IsXMLSchemaDateTime() {
			p.Created = pp.Get()
		}
	}

	p.Type = s.GetTypeName()
	return p
}
