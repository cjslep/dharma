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
	"net/url"

	"github.com/go-fed/activity/streams/vocab"
)

type Authorable interface {
	GetActivityStreamsAttributedTo() vocab.ActivityStreamsAttributedToProperty
}

func ToAuthors(a Authorable) []*url.URL {
	atp := a.GetActivityStreamsAttributedTo()
	if atp == nil {
		return nil
	}
	authors := make([]*url.URL, 0, atp.Len())
	for iter := atp.Begin(); iter != atp.End(); iter = iter.Next() {
		if iter.IsIRI() {
			authors = append(authors, iter.GetIRI())
		} else if iter.IsActivityStreamsPerson() {
			pidp := iter.GetActivityStreamsPerson().GetJSONLDId()
			if pidp != nil {
				authors = append(authors, pidp.GetIRI())
			}
		}
	}
	return authors
}
