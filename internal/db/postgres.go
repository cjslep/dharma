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

	"github.com/go-fed/apcore/app"
)

type postgres struct {
	schema string // with trailing "."
}

func newPostgres(schema string) postgres {
	if len(schema) > 0 {
		schema = schema + "."
	}
	return postgres{schema}
}

func CreateTables(c context.Context, db app.Database, schema string) error {
	p := newPostgres(schema)
	tx := db.Begin()
	tx.Exec(p.CreateEvePublicKeysTableV0())
	tx.Exec(p.CreateEveTokensTableV0())
	tx.Exec(p.CreateApplicationStateTableV0())
	tx.Exec(p.CreateUserSupplementTableV0())
	return tx.Do(c)
}

// EVE Online Public Keys Table

func (p postgres) CreateEvePublicKeysTableV0() string {
	return `
CREATE TABLE IF NOT EXISTS ` + p.schema + `dharma_eve_public_keys
(
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  create_time timestamp with time zone DEFAULT current_timestamp,
  hash text UNIQUE NOT NULL,
  keys jsonb NOT NULL
);`
}

func (p postgres) AddEvePublicKey() string {
	return `INSERT INTO ` + p.schema + `dharma_eve_public_keys
(hash, keys)
VALUES
(SHA512($2), $1)
ON CONFLICT (hash) DO NOTHING;`
}

func (p postgres) GetLatestEvePublicKey() string {
	return `SELECT keys FROM ` + p.schema + `dharma_eve_public_keys
ORDER BY create_time DESC NULLS LAST
LIMIT 1;`
}

// EVE Online Tokens Table

func (p postgres) CreateEveTokensTableV0() string {
	return `
CREATE TABLE IF NOT EXISTS ` + p.schema + `dharma_eve_tokens
(
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid REFERENCES ` + p.schema + `users(id) ON DELETE CASCADE NOT NULL,
  character_id integer UNIQUE NOT NULL,
  status text NOT NULL,
  tokens jsonb NOT NULL
);`
}

func (p postgres) SetEveToken() string {
	return `INSERT INTO ` + p.schema + `dharma_eve_tokens
(user_id, character_id, status, tokens)
VALUES
($1, $2, $4, $3)
ON CONFLICT (character_id) DO UPDATE
SET tokens = EXCLUDED.tokens;`
}

func (p postgres) GetEveToken() string {
	return `SELECT tokens FROM ` + p.schema + `dharma_eve_tokens
WHERE character_id = $1;`
}

func (p postgres) GetExpiringEveTokensWithin() string {
	return `SELECT user_id, tokens FROM ` + p.schema + `dharma_eve_tokens
WHERE (tokens->access_expires)::timestamp < current_timestamp + interval '$1'`
}

func (p postgres) MarkAllTokensWithState() string {
	return `UPDATE ` + p.schema + `dharma_eve_tokens
(status) VALUES ($1)`
}

func (p postgres) GetEveCharactersForUser() string {
	return `SELECT character_id, status FROM ` + p.schema + `dharma_eve_tokens
WHERE user_id = $1;`
}

func (p postgres) HasCharacterForUser() string {
	return `EXISTS(SELECT id FROM ` + p.schema + `dharma_eve_tokens
WHERE user_id = $1 AND character_id = $2);`
}

func (p postgres) GetTokenStatusForCharacter() string {
	return `SELECT status FROM ` + p.schema + `dharma_eve_tokens
WHERE character_id = $1);`
}

// Application State Table

func (p postgres) CreateApplicationStateTableV0() string {
	return `
CREATE TABLE IF NOT EXISTS ` + p.schema + `dharma_application_state
(
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  kind text UNIQUE NOT NULL,
  value text NOT NULL
);`
}

func (p postgres) SetApplicationStateKV() string {
	return `INSERT INTO ` + p.schema + `dharma_application_state
(kind, value)
VALUES
($1, $2)
ON CONFLICT (kind) DO UPDATE
SET value = EXCLUDED.value;`
}

func (p postgres) GetApplicationStateKV() string {
	return `SELECT value FROM ` + p.schema + `dharma_application_state
WHERE kind = $1;`
}

func (p postgres) GetApplicationStateValueAndPrefix() string {
	return `SELECT kind FROM ` + p.schema + `dharma_application_state
WHERE kind LIKE CONCAT($1, "%") AND value = $2`
}

// User Supplementary Data Table

func (p postgres) CreateUserSupplementTableV0() string {
	return `
CREATE TABLE IF NOT EXISTS ` + p.schema + `dharma_user_supplement
(
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid REFERENCES ` + p.schema + `users(id) ON DELETE CASCADE NOT NULL,
  validation_state text NOT NULL,
  verification_token text NOT NULL
);`
}

func (p postgres) CreateUserSupplement() string {
	return `
INSERT INTO ` + p.schema + `dharma_user_supplement
(user_id, verification_token, validation_state)
VALUES
($1, $2, $3);`
}

func (p postgres) UpdateUserSupplement() string {
	return `UPDATE ` + p.schema + `dharma_user_supplement
SET validation_state = $2 WHERE user_id = $1;`
}

func (p postgres) UpdateUserVerified() string {
	return `UPDATE ` + p.schema + `dharma_user_supplement
SET validation_state = $2 WHERE verification_token = $1;`
}

func (p postgres) GetUserSupplement() string {
	return `SELECT validation_state FROM ` + p.schema + `dharma_user_supplement
WHERE user_id = $1;`
}
