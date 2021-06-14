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
	"net/url"
	"strings"
)

type OAuth2Client struct {
	RedirectURI string
	ClientID    string
}

func (o *OAuth2Client) GetURL(state string, scopes []string) *url.URL {
	u := &url.URL{
		Scheme: "https",
		Host:   "login.eveonline.com",
		Path:   "/v2/oauth/authorize/",
	}
	v := url.Values{}
	v.Add("response_type", "code")
	v.Add("redirect_uri", o.RedirectURI)
	v.Add("client_id", o.ClientID)
	v.Add("scope", strings.Join(scopes, ","))
	v.Add("state", state)
	u.RawQuery = v.Encode()
	return u
}
