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
	"sort"

	"github.com/cjslep/dharma/internal/data"
	"github.com/cjslep/dharma/internal/db"
	"github.com/go-fed/apcore/util"
	"golang.org/x/text/language"
)

type Tags struct {
	DB *db.DB
}

type LatestTag struct {
	T data.Tag
	S []data.Snippet
}

func (l *LatestTag) sort() {
	sort.Sort(data.LatestSnippets(l.S))
}

// GetLatestSnippets obtains the latest snippets
func (t *Tags) GetLatestSnippets(ctx util.Context, display []data.Tag, n, length, maxHtmlDepth int, preferLang language.Tag) (map[string]*LatestTag, error) {
	l, err := t.DB.FetchLatestPublicTags(ctx, display, n)
	if err != nil {
		return nil, err
	}

	// Ensure default data is available
	lt := make(map[string]*LatestTag, len(display))
	for _, tag := range display {
		lt[tag.ID] = &LatestTag{T: tag}
	}

	// Populate with data
	for _, lpt := range l {
		snip := data.ToSnippet(lpt.T, length, maxHtmlDepth, preferLang)
		tags := data.ToTags(lpt.T)
		for _, tag := range tags {
			v := lt[tag.ID]
			if v == nil {
				v = &LatestTag{}
				lt[tag.ID] = v
			}
			v.S = append(v.S, snip)
		}
	}

	// Sort each category
	for _, v := range lt {
		v.sort()
	}
	return lt, nil
}

// GetThreadsPreviewsForTag obtains the latest threads for a given tag
func (t *Tags) GetThreadPreviewsForTag(ctx util.Context, g data.Tag, n, maxDepth, page int, preferLang language.Tag) ([]data.ThreadPreview, error) {
	h, err := t.DB.FetchMostRecentlyUpdatedThreads(ctx, g, n, page)
	if err != nil {
		return nil, err
	}
	tps := make([]data.ThreadPreview, 0, len(h))
	for _, r := range h {
		tps = append(tps, data.ToThreadPreview(r.First, r.MostRecent, n, maxDepth, preferLang))
	}
	sort.Sort(data.RecentPreviews(tps))
	return tps, nil
}
