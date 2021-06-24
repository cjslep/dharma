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

// SoftMapMerge takes m[k], converts to map[string]interface{}, and merges val
// into m[k], overwriting if necessary.
//
// If m[k] does not exist, it is created.
//
// If it cannot do the merge, silently fail.
func SoftMapMergeInto(m map[string]interface{}, key string, val map[string]interface{}) {
	i, ok := m[key]
	if !ok {
		i = make(map[string]interface{})
		m[key] = i
	}
	if mk, ok := i.(map[string]interface{}); ok {
		m[key] = MapMerge(mk, val)
	}
}

// MapMerge merges 'from' into 'to', overwriting if necessary.
func MapMerge(to, from map[string]interface{}) map[string]interface{} {
	for k, v := range from {
		to[k] = v
	}
	return to
}
