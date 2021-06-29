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

type LatestSnippets []Snippet

func (l LatestSnippets) Len() int {
	return len(l)
}

func (l LatestSnippets) Less(i, j int) bool {
	return l[i].Created.After(l[j].Created)
}

func (l LatestSnippets) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type RecentPreviews []ThreadPreview

func (r RecentPreviews) Len() int {
	return len(r)
}

func (r RecentPreviews) Less(i, j int) bool {
	return r[i].Last.Created.After(r[j].Last.Created)
}

func (r RecentPreviews) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
