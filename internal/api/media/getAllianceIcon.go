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

package media

import (
	"fmt"
	"image"
	"net/http"
	"net/url"

	"github.com/cjslep/dharma/internal/api/paths"
	"github.com/cjslep/dharma/internal/services"
)

func (m *Media) getAllianceIcon(w http.ResponseWriter, r *http.Request) {
	m.handleEveImage(w, r, func(aID int32, query url.Values) (string, image.Image, error) {
		var size services.AllianceIconSize
		qsize := query.Get(paths.EveMediaSizeParam)
		switch qsize {
		default:
			fallthrough
		case "64":
			size = services.A64
		case "128":
			size = services.A128
		}

		img, err := m.C.Media.GetAllianceIcon(r.Context(), aID, size)
		name := fmt.Sprintf("alliance_%d", aID)
		return name, img, err
	})
}
