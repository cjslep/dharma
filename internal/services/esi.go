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
	"github.com/cjslep/dharma/esi"
	"github.com/cjslep/dharma/internal/db"
)

type ESI struct {
	DB *db.DB
}

func (e *ESI) SetEvePublicKeys(o *esi.OAuthKeysMetadata) error {
	// TODO
	return e.DB.SetEvePublicKeys(o)
}

func (e *ESI) GetEvePublicKeys() (*esi.OAuthKeysMetadata, error) {
	// TODO
	return e.DB.GetEvePublicKeys()
}

func (e *ESI) SetEveTokens(t *esi.Tokens) error {
	// TODO
	return e.DB.SetEveTokens(t)
}

func (e *ESI) GetEveTokens() (*esi.Tokens, error) {
	// TODO
	return e.DB.GetEveTokens()
}
