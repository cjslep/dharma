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

package services

import (
	"context"
	"strings"
	"time"

	"github.com/cjslep/dharma/esi"
	"github.com/cjslep/dharma/internal/async"
	"github.com/cjslep/dharma/internal/db"
	"github.com/go-fed/apcore/util"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
)

type ESI struct {
	DB        *db.DB
	OAC       *esi.OAuth2Client
	L         *zerolog.Logger
	ESIClient *esi.Client
}

func (e *ESI) SetEvePublicKeys(c util.Context, o *esi.OAuthKeysMetadata) error {
	return e.DB.SetEvePublicKeys(c, o)
}

func (e *ESI) GetEvePublicKeys(c util.Context) (*esi.OAuthKeysMetadata, error) {
	return e.DB.GetEvePublicKeys(c)
}

func (e *ESI) SetEveTokens(c util.Context, userID string, t *esi.Tokens) error {
	return e.DB.SetEveTokens(c, userID, t)
}

func (e *ESI) GetEveTokens(c util.Context) (*esi.Tokens, error) {
	return e.DB.GetEveTokens(c)
}

func (e *ESI) GoPeriodicallyRefreshAllTokens(m *async.Messenger) {
	m.Periodically(time.Hour, e.refreshAllTokens, e.L) // TODO: Configurable time
}

func (e *ESI) refreshAllTokens(c context.Context) error {
	// TODO
	return nil
}

func (e *ESI) SearchCorporations(c context.Context, query string, lang language.Tag) ([]*esi.Corporation, error) {
	s := strings.TrimSpace(query)
	if len(s) <= 3 {
		return nil, errors.New("query is too short to search")
	}
	return e.ESIClient.SearchCorp(c, s, lang)
}

func (e *ESI) GetCharactersForUser(c util.Context, userID string) ([]*esi.Character, error) {
	ids, err := e.DB.GetEveCharactersForUser(c, userID)
	if err != nil {
		return nil, err
	}
	return e.ESIClient.Characters(c, ids)
}

func (e *ESI) HasCharacterForUser(c util.Context, userID string, charID int32) (bool, error) {
	return e.DB.HasCharacterForUser(c, userID, charID)
}
