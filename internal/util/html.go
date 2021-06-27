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

package util

import (
	"bytes"

	"golang.org/x/net/html"
)

func FirstNHTML(text string, n, maxDepth int) string {
	in := bytes.NewBuffer([]byte(text))
	var out bytes.Buffer
	z := html.NewTokenizer(in)
	depthTags := make([]string, 0, maxDepth)
	l := 0
	// 1. Extract a sizeable text to display
	for l < n {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// Malformed
			break
		case html.TextToken:
			l += len(z.Text())
			if l > (n + n/2) {
				// Too long text
				break
			}
		case html.StartTagToken:
			if len(depthTags) == maxDepth {
				// Too many nested HTML
				break
			}
			tag, _ := z.TagName()
			depthTags = append(depthTags, string(tag))
		case html.EndTagToken:
			depthTags = depthTags[:len(depthTags)-1]
		}
		out.Write(z.Raw())
	}
	// 2. Close the rest of the tags
	for i := len(depthTags) - 1; i >= 0; i-- {
		out.WriteString("</")
		out.WriteString(depthTags[i])
		out.WriteString(">")
	}
	return out.String()
}
