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

package esi

import (
	"image"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// FetchEveOnlineImage fetches an image from EvE Online.
//
// It is able to decode JPEG, PNG, and GIF images.
//
// The expires determines the ability to cache for our server. It is possible
// it is the zero Time, in which case we could not determine the cacheability.
func FetchEveOnlineImage(u *url.URL) (i image.Image, expires time.Time, err error) {
	var resp *http.Response
	resp, err = http.Get(u.String())
	if err != nil {
		return
	}
	// Be conservative by using the basis of our cache-expiry time to be
	// later than the true expiry time, which would be during the GET call.
	now := time.Now()
	// Produce the image
	i /*image type=*/, _, err = image.Decode(resp.Body)
	if err != nil {
		return
	}
	// Determine the expiration time
	sage := resp.Header.Get("Age")
	cc := resp.Header.Get("Cache-Control")
	if len(sage) == 0 || len(cc) == 0 || !strings.HasPrefix(cc, "max-age=") {
		return
	}
	var age, maxAge int
	if age, err = strconv.Atoi(sage); err != nil {
		return
	}
	if maxAge, err = strconv.Atoi(strings.TrimPrefix(cc, "max-age=")); err != nil {
		return
	}
	remaining := maxAge - age
	expires = now.Add(time.Duration(remaining) * time.Second)
	return
}
