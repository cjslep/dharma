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

type InReplyToable interface {
	GetActivityStreamsInReplyTo() vocab.ActivityStreamsInReplyToProperty
}

func ToInReplyTo(i InReplyToable) []*url.URL {
	irtp := i.GetActivityStreamsInReplyTo()
	if irtp == nil {
		return nil
	}
	irt := make([]*url.URL, 0, irtp.Len())
	for iter := irtp.Begin(); iter != irtp.End(); iter = iter.Next() {
		if iter.IsIRI() {
			irt = append(irt, iter.GetIRI())
		}
	}
	return irt
}
