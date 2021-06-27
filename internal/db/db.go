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

package db

import (
	"context"
	"time"

	"github.com/cjslep/dharma/esi"
	"github.com/cjslep/dharma/internal/data"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/go-fed/apcore/app"
)

type DB struct {
	db app.Database
}

func New(db app.Database) *DB {
	return &DB{
		db: db,
	}
}

func (d *DB) SetEvePublicKeys(c context.Context, o *esi.OAuthKeysMetadata) error {
	// TODO
	return nil
}

func (d *DB) GetEvePublicKeys(c context.Context) (*esi.OAuthKeysMetadata, error) {
	// TODO
	return nil, nil
}

func (d *DB) SetEveTokens(c context.Context, t *esi.Tokens) error {
	// TODO
	return nil
}

func (d *DB) GetEveTokens(c context.Context) (*esi.Tokens, error) {
	// TODO
	return nil, nil
}

type LatestPublicTagsResult struct {
	T        vocab.Type
	Received time.Time
}

func (d *DB) FetchLatestPublicTags(c context.Context, display []data.Tag, n int) ([]LatestPublicTagsResult, error) {
	// TODO
	return nil, nil
}
