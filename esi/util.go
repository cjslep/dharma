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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cjslep/dharma/internal/features"
)

type OAuth2Client struct {
	RedirectURI string
	ClientID    string
	Secret      string
	Client      *http.Client
}

func (o *OAuth2Client) GetURL(state string, scopes []features.Scope) *url.URL {
	u := &url.URL{
		Scheme: "https",
		Host:   "login.eveonline.com",
		Path:   "/v2/oauth/authorize/",
	}
	s := make([]string, len(scopes))
	for i := range scopes {
		s[i] = string(scopes[i])
	}
	v := url.Values{}
	v.Add("response_type", "code")
	v.Add("redirect_uri", o.RedirectURI)
	v.Add("client_id", o.ClientID)
	v.Add("scope", strings.Join(s, ","))
	v.Add("state", state)
	u.RawQuery = v.Encode()
	return u
}

func (o *OAuth2Client) GetAuthorization(code string) (*JWTResponse, error) {
	// Issue request
	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	req, err := http.NewRequest(
		"POST",
		"https://login.eveonline.com/v2/oauth/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Host = "login.eveonline.com"
	rawAuth := fmt.Sprintf("%s:%s", o.ClientID, o.Secret)
	auth := base64.URLEncoding.EncodeToString([]byte(rawAuth))
	req.Header = map[string][]string{
		"Content-Type":  {"application/x-www-form-urlencoded"},
		"Authorization": {fmt.Sprintf("Basic %s", auth)},
	}
	resp, err := o.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Process payload body
	respb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	jwt := &JWTResponse{}
	err = json.Unmarshal(respb, jwt)
	if err != nil {
		return nil, err
	}
	return jwt, nil
}

func (o *OAuth2Client) GetRefresh(refresh string) (*JWTResponse, error) {
	// Issue request
	data := url.Values{}
	data.Add("grant_type", "refresh_token")
	data.Add("refresh_token", refresh)
	req, err := http.NewRequest(
		"POST",
		"https://login.eveonline.com/v2/oauth/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Host = "login.eveonline.com"
	rawAuth := fmt.Sprintf("%s:%s", o.ClientID, o.Secret)
	auth := base64.URLEncoding.EncodeToString([]byte(rawAuth))
	req.Header = map[string][]string{
		"Content-Type":  {"application/x-www-form-urlencoded"},
		"Authorization": {fmt.Sprintf("Basic %s", auth)},
	}
	resp, err := o.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Process payload body
	respb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	jwt := &JWTResponse{}
	err = json.Unmarshal(respb, jwt)
	if err != nil {
		return nil, err
	}
	return jwt, nil
}
