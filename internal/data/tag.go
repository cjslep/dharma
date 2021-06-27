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
	"github.com/go-fed/activity/streams/vocab"
)

var (
	// These tags are hardcoded for now, as custom tagging will require a
	// lot of work and special logic for "other" or "uncategorized"

	// BEGIN IF THESE ARE CHANGED: Change forum/home.tmpl
	// Corp
	Announce = Tag{"announce"}
	Events   = Tag{"events"}
	Discuss  = Tag{"discuss"}

	// Game
	Fleet    = Tag{"fleet"}
	Industry = Tag{"industry"}
	Market   = Tag{"market"}
	PVP      = Tag{"pvp"}
	PVE      = Tag{"pve"}

	// Social
	Relations = Tag{"relations"}
	Intel     = Tag{"intel"}
	Justice   = Tag{"justice"}
	QNA       = Tag{"qna"}

	// Other
	OffTopic      = Tag{"offtopic"}
	Uncategorized = Tag{"uncategorized"}
	// END IF THESE ARE CHANGED: Change forum/home.tmpl
)

var AllTags = []Tag{
	Announce,
	Events,
	Discuss,
	Fleet,
	Industry,
	Market,
	PVP,
	PVE,
	Intel,
	Relations,
	QNA,
	Justice,
	OffTopic,
	Uncategorized,
}

type Tag struct {
	ID string // Is also the user-displayable string
}

type tagable interface {
	GetActivityStreamsTag() vocab.ActivityStreamsTagProperty
}

func ToTags(t vocab.Type) []Tag {
	s, ok := t.(tagable)
	if !ok {
		return nil
	}

	// var g Tag
	tp := s.GetActivityStreamsTag()
	if tp != nil {
		// TODO: Our own category tag type?
	}

	return nil
}
