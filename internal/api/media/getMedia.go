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
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-fed/apcore/app"
	"github.com/gorilla/mux"
)

// TODO: Support Privileges, and determine whether the request is privileged
// enough to fetch the media.
func (m *Media) getMedia(w http.ResponseWriter, r *http.Request, k app.Session) {
	id := mux.Vars(r)["id"]
	name, contentType, b, err := m.C.Media.GetMedia(r.Context(), id)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		m.C.L.Error().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	read := bytes.NewReader(b)

	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", name))
	w.Header().Set("Content-Type", contentType)
	// TODO: fetch the creation time to use as the modtime
	modtime := time.Time{}
	http.ServeContent(w, r, "", modtime, read)
}
