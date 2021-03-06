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

package paths

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/cjslep/dharma/internal/features"
	"golang.org/x/text/language"
)

const (
	TokenQueryParam       = "t"
	VerifyTYParam         = "ty"
	VerifySuccessParam    = "s"
	RescopeQueryParam     = "rescope"
	EveMediaSizeParam     = "size"
	VerifyPath            = "/account/verify"
	ESIAuthPath           = "/esi/auth"
	AccountCharactersPath = "/account/characters"
)

func LocalizedRoot(lang language.Tag) *url.URL {
	u := &url.URL{
		Path: fmt.Sprintf("/%s", lang),
	}
	return u
}

func tokenizeVerifyPath(scheme, host, token string, lang language.Tag) string {
	var v url.Values
	v.Set(TokenQueryParam, token)
	u := &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     fmt.Sprintf("/%s%s", lang, VerifyPath),
		RawQuery: v.Encode(),
	}
	return u.String()
}

func TokenizeVerifyPath(host, token string, lang language.Tag) string {
	return tokenizeVerifyPath("https", host, token, lang)
}

func TokenizeVerifyPathHTTP(host, token string, lang language.Tag) string {
	return tokenizeVerifyPath("http", host, token, lang)
}

func IsVerifyPath(path string) bool {
	return strings.HasSuffix(path, VerifyPath)
}

func GetPleaseVerifyURLWithTY(lang language.Tag) *url.URL {
	return getPleaseVerifyURL(lang, true, false)
}

func GetPleaseVerifyURLWithSuccess(lang language.Tag) *url.URL {
	return getPleaseVerifyURL(lang, false, true)
}

func getPleaseVerifyURL(lang language.Tag, showThanksForRegistering, showVerifySuccess bool) *url.URL {
	u := &url.URL{}
	u.Path = fmt.Sprintf("/%s%s", lang, VerifyPath)
	v := url.Values{}
	if showThanksForRegistering {
		v.Add(VerifyTYParam, "true")
	}
	if showVerifySuccess {
		v.Add(VerifySuccessParam, "true")
	}
	u.RawQuery = v.Encode()
	return u
}

func GetPostESIAuthPath(lang language.Tag, l features.List) *url.URL {
	v := url.Values{}
	for _, id := range l.IDs() {
		v.Add("features", id)
	}
	u := &url.URL{
		Path:     fmt.Sprintf("/%s%s", lang, ESIAuthPath),
		RawQuery: v.Encode(),
	}
	return u
}

func GetESIAuthPathRescope(lang language.Tag) *url.URL {
	v := url.Values{}
	v.Add(RescopeQueryParam, "true")
	u := &url.URL{
		Path:     fmt.Sprintf("/%s%s", lang, ESIAuthPath),
		RawQuery: v.Encode(),
	}
	return u
}

func GetCharacterSelection(lang language.Tag) *url.URL {
	u := &url.URL{
		Path: fmt.Sprintf("/%s%s", lang, AccountCharactersPath),
	}
	return u
}
