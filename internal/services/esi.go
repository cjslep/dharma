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
	"time"

	"github.com/cjslep/dharma/esi"
	"github.com/cjslep/dharma/internal/async"
	"github.com/cjslep/dharma/internal/db"
	"github.com/rs/zerolog"
)

type ESI struct {
	DB  *db.DB
	OAC *esi.OAuth2Client
	L   *zerolog.Logger
}

func (e *ESI) SetEvePublicKeys(c context.Context, o *esi.OAuthKeysMetadata) error {
	// TODO
	return e.DB.SetEvePublicKeys(c, o)
}

func (e *ESI) GetEvePublicKeys(c context.Context) (*esi.OAuthKeysMetadata, error) {
	// TODO
	return e.DB.GetEvePublicKeys(c)
}

func (e *ESI) SetEveTokens(c context.Context, t *esi.Tokens) error {
	// TODO
	return e.DB.SetEveTokens(c, t)
}

func (e *ESI) GetEveTokens(c context.Context) (*esi.Tokens, error) {
	// TODO
	return e.DB.GetEveTokens(c)
}

func (e *ESI) GoPeriodicallyRefreshAllTokens(m *async.Messenger) {
	m.Periodically(time.Hour, e.refreshAllTokens, e.L) // TODO: Configurable time
}

func (e *ESI) refreshAllTokens(c context.Context) error {
	// TODO
	return nil
}
