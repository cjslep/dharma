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

package extract

import (
	"github.com/cjslep/dharma/internal/util"
	"github.com/go-fed/activity/streams/vocab"
	"golang.org/x/text/language"
)

type Contentable interface {
	GetActivityStreamsContent() vocab.ActivityStreamsContentProperty
}

func ToPreviewContent(c Contentable, n, maxDepth int, preferLang language.Tag) string {
	cp := c.GetActivityStreamsContent()
	if cp == nil {
		return ""
	}
	for iter := cp.Begin(); iter != cp.End(); iter = iter.Next() {
		// TODO: mime type HTML and Markdown
		if iter.IsXMLSchemaString() {
			return util.FirstNHTML(iter.GetXMLSchemaString(), n, maxDepth)
		} else if iter.IsRDFLangString() {
			if iter.HasLanguage(preferLang.String()) {
				return util.FirstNHTML(iter.GetLanguage(preferLang.String()), n, maxDepth)
			} else if iter.HasLanguage("en") {
				return util.FirstNHTML(iter.GetLanguage("en"), n, maxDepth)
			}
		}
	}
	return ""
}
