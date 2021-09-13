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
	"fmt"
	"image"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/cjslep/dharma/internal/util"
	"github.com/gorilla/mux"
)

// fetchImageFn is given a character ID and the request query parameters, and
// should return an image of the appropriate kind and name for its
// implementation.
type fetchImageFn func(id int32, query url.Values) (string, image.Image, error)

// handleEveImage abstracts the steps required for a handler that:
// 1. Has an {id} in its request path
// 2. Can fetch an image from ESI / CCP Static serving using this ID
// 3. Can determine the size of the requested image from the query parameters
// 4. Wishes to encode the resulting image in a common encoding
func (m *Media) handleEveImage(w http.ResponseWriter, r *http.Request, fifn fetchImageFn) {
	id := mux.Vars(r)["id"]
	pid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		m.C.L.Error().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pid32 := int32(pid)
	q := r.URL.Query()

	name, img, err := fifn(pid32, q)
	if err != nil {
		m.C.L.Error().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := util.EncodeJPEG(img)
	if err != nil {
		m.C.L.Error().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	read := bytes.NewReader(b)

	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.jpeg\"", name))
	w.Header().Set("Content-Type", "image/jpeg")
	modtime := time.Time{}
	http.ServeContent(w, r, "", modtime, read)
}
