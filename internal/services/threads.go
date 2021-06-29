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
	"context"
	"sort"

	"github.com/cjslep/dharma/internal/data"
	"github.com/cjslep/dharma/internal/db"
	"golang.org/x/text/language"
)

type Threads struct {
	DB *db.DB
}

func (t *Threads) GetPosts(ctx context.Context, id string, n, page int, preferLang language.Tag) ([]data.Post, error) {
	m, err := t.DB.FetchPaginatedMessagesInThread(ctx, id, n, page)
	if err != nil {
		return nil, err
	}
	ts := make([]data.Post, 0, len(m.Messages))
	for _, r := range m.Messages {
		ts = append(ts, data.ToPost(r, preferLang))
	}
	sort.Sort(data.ChronologicalPosts(ts))
	return ts, nil
}
