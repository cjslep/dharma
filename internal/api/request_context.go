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

package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cjslep/dharma/internal/sessions"
	"github.com/go-fed/apcore/app"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

const (
	pathContextValue         = "pcv"
	sessionContextValue      = "scv"
	languageTagsContextValue = "ltcv"
)

type RequestContext struct {
	context.Context
}

func From(c context.Context) *RequestContext {
	return &RequestContext{c}
}

func (r *RequestContext) Update(req *http.Request) *http.Request {
	return req.WithContext(r.Context)
}

func (r *RequestContext) WithSession(k app.Session) {
	r.Context = context.WithValue(r.Context, sessionContextValue, k)
}

func (r *RequestContext) Session() (app.Session, error) {
	v := r.Context.Value(sessionContextValue)
	if v == nil {
		return nil, errors.New("no session in request context")
	} else if k, ok := v.(app.Session); !ok {
		return nil, errors.Errorf("request context session is not app.Session: %T", v)
	} else {
		return k, nil
	}
}

func (r *RequestContext) WithLanguageTags(l ...language.Tag) {
	r.Context = context.WithValue(r.Context, languageTagsContextValue, l)
}

func (r *RequestContext) LanguageTags() ([]language.Tag, error) {
	v := r.Context.Value(languageTagsContextValue)
	if v == nil {
		return nil, errors.New("no language tags in request context")
	} else if k, ok := v.([]language.Tag); !ok {
		return nil, errors.Errorf("request context language tags are not []language.Tag: %T", v)
	} else {
		return k, nil
	}
}

func (r *RequestContext) WithPath(p string) {
	r.Context = context.WithValue(r.Context, pathContextValue, p)
}

func (r *RequestContext) Path() (string, error) {
	v := r.Context.Value(pathContextValue)
	if v == nil {
		return "", errors.New("no path in request context")
	} else if k, ok := v.(string); !ok {
		return "", errors.Errorf("request context path is not string: %T", v)
	} else {
		return k, nil
	}
}

func (r *RequestContext) navData(signedIn bool, tag language.Tag, charID int32) map[string]interface{} {
	m := map[string]interface{}{
		"signedIn": signedIn,
		"paths": map[string]interface{}{
			"register":        fmt.Sprintf("/%s/account/register", tag),
			"login":           fmt.Sprintf("/%s/login", tag),
			"logout":          fmt.Sprintf("/%s/logout", tag),
			"changeCharacter": fmt.Sprintf("/%s/account/characters", tag),
			"profile":         fmt.Sprintf("/%s/account/profile", tag),
			"settings":        fmt.Sprintf("/%s/account/settings", tag),
			"forum":           fmt.Sprintf("/%s/forum", tag),
			"killboard":       fmt.Sprintf("/%s/killboard", tag),
			"calendar":        fmt.Sprintf("/%s/calendar", tag),
		},
		"localizePath": func(s string) (string, error) {
			p, err := r.Path()
			if err != nil {
				return "", err
			}
			k := strings.Split(p, "/")
			if len(k) < 2 {
				return s, nil
			}
			k[1] = s
			return strings.Join(k, "/"), nil
		},
	}
	if charID != 0 {
		m["characterID"] = charID
	}
	return m
}

func (r *RequestContext) RenderNavData() map[string]interface{} {
	// Determine signed-in state
	k, err := r.Session()
	signedIn := false
	var charID int32
	if err == nil {
		// Determine if signed in
		_, isSignedInErr := k.UserID()
		signedIn = isSignedInErr == nil
		// Determine current character
		charID = sessions.GetCharacterSelected(k)
	}

	// Obtain a language
	ts, err := r.LanguageTags()
	tag := language.English
	if err == nil {
		tag = ts[0]
	}

	return r.navData(signedIn, tag, charID)
}
