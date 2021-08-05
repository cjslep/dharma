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
	dutil "github.com/cjslep/dharma/internal/util"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
)

type ESI struct {
	DB               *db.DB
	OAC              *esi.OAuth2Client
	L                *zerolog.Logger
	ESIClient        *esi.Client
	PeriodicRefresh  time.Duration
	PeriodicKeyFetch time.Duration
}

func (e *ESI) GoPeriodicallyFetchEvePublicKeys(m *async.Messenger) {
	m.NowAndPeriodically(e.PeriodicKeyFetch, e.fetchEveOnlineKeys, e.L)
}

func (e *ESI) fetchEveOnlineKeys(c context.Context) error {
	keys, err := esi.FetchEveOnlineKeys()
	if err != nil {
		return err
	}
	return e.DB.SetEvePublicKeys(c, keys)
}

func (e *ESI) GetEvePublicKeys(c context.Context) (*esi.OAuthKeysMetadata, error) {
	return e.DB.GetEvePublicKeys(c)
}

func (e *ESI) SetEveTokens(c context.Context, userID string, t *esi.Tokens) error {
	return e.DB.SetEveTokens(c, userID, t)
}

func (e *ESI) GetEveToken(c context.Context, charID int32) (*esi.Tokens, error) {
	return e.DB.GetEveToken(c, charID)
}

func (e *ESI) GoPeriodicallyRefreshAllTokens(m *async.Messenger) {
	m.Periodically(e.PeriodicRefresh, e.refreshAllTokens, e.L)
}

func (e *ESI) refreshAllTokens(c context.Context) error {
	uts, err := e.DB.GetExpiringEveTokensWithin(c, e.PeriodicRefresh)
	if err != nil {
		return err
	}
	errs := make([]error, 0, len(uts))
	for i, ut := range uts {
		errs[i] = e.refreshToken(c, ut.UserID, ut.T.Refresh)
	}
	return dutil.ToErrors(errs)
}

func (e *ESI) refreshToken(c context.Context, userID string, refresh string) error {
	// Do the refresh
	jwt, err := e.OAC.GetRefresh(refresh)
	if err != nil {
		return err
	}

	// Verify the authenticity of the new authorization
	ek, err := e.GetEvePublicKeys(c)
	if err != nil {
		return err
	}
	jwtk := ek.JWTKey()
	if jwtk == nil {
		return errors.New("could not find EVE public jwt key")
	}
	claims, err := jwtk.ValidateToken([]byte(jwt.AccessToken))
	if err != nil {
		return err
	}
	err = esi.ValidateEveClaims(claims)
	if err != nil {
		return err
	}

	// Construct our internal representation of a validated token, and
	// store it.
	tokens, err := esi.NewTokens(jwt, claims)
	if err != nil {
		return err
	}
	err = e.SetEveTokens(c, userID, tokens)
	if err != nil {
		return err
	}
	return nil
}

func (e *ESI) SearchCorporations(c context.Context, query string, lang language.Tag) ([]*esi.Corporation, error) {
	s := strings.TrimSpace(query)
	if len(s) <= 3 {
		return nil, errors.New("query is too short to search")
	}
	return e.ESIClient.SearchCorp(c, s, lang)
}

func (e *ESI) GetCharactersForUser(c context.Context, userID string) ([]*esi.Character, error) {
	ids, err := e.DB.GetEveCharactersForUser(c, userID)
	if err != nil {
		return nil, err
	}
	return e.ESIClient.Characters(c, ids)
}

func (e *ESI) HasCharacterForUser(c context.Context, userID string, charID int32) (bool, error) {
	return e.DB.HasCharacterForUser(c, userID, charID)
}
