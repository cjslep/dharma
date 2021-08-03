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
	"fmt"
	"strings"
)

var _ error = &manyErrors{}

type manyErrors []error

func ToErrors(nilable []error) error {
	notNilErrs := make([]error, 0, len(nilable))
	for _, err := range nilable {
		if err != nil {
			notNilErrs = append(notNilErrs, err)
		}
	}
	if len(notNilErrs) == 0 {
		return nil
	}
	e := manyErrors(notNilErrs)
	return &e
}

func (e *manyErrors) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d errors occurred:", len(*e))
	for _, err := range *e {
		fmt.Fprintf(&b, "\n> %s", err)
	}
	return b.String()
}
