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
	"time"

	"github.com/cjslep/dharma/esi"
	"github.com/cjslep/dharma/internal/data"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/go-fed/apcore/app"
	"github.com/go-fed/apcore/util"
)

const (
	// Email validation states
	kUnvalidatedState             = "not_validated"
	kSentValidationChallengeState = "challenge_sent"
	kValidatedState               = "validated"
)

const (
	// Application state Key-Val keys
	kAuthoritativeCharacterKey = "authoritative_character"
	kCorporationManagedKey     = "corporation_managed"
	kAllianceAssociationKey    = "alliance_association"
	kExecutorCorporationKey    = "executor_corporation"
)

type DB struct {
	db    app.Database
	f     app.Framework
	pg    postgres
	cache *cache
}

func New(db app.Database, f app.Framework, schema string) *DB {
	return &DB{
		db:    db,
		f:     f,
		pg:    newPostgres(schema),
		cache: newCache(),
	}
}

func (d *DB) SetEvePublicKeys(c util.Context, o *esi.OAuthKeysMetadata) error {
	txb := d.db.Begin()
	txb.Exec(d.pg.AddEvePublicKey(), o)
	return txb.Do(c)
}

func (d *DB) GetEvePublicKeys(c util.Context) (*esi.OAuthKeysMetadata, error) {
	// TODO: See if in-memory cache is needed for per-request latency
	o := &esi.OAuthKeysMetadata{}
	txb := d.db.Begin()
	txb.QueryOneRow(d.pg.GetLatestEvePublicKey(), func(r app.SingleRow) error {
		return r.Scan(o)
	})
	return o, txb.Do(c)
}

func (d *DB) SetEveTokens(c util.Context, t *esi.Tokens) error {
	txb := d.db.Begin()
	txb.ExecOneRow(d.pg.SetEveToken(), t.CID, t)
	return txb.Do(c)
}

func (d *DB) GetEveTokens(c util.Context) (*esi.Tokens, error) {
	// TODO: Be more intelligent about refreshing keys and detecting cacheability
	t := &esi.Tokens{}
	txb := d.db.Begin()
	txb.QueryOneRow(d.pg.GetEveToken(), func(r app.SingleRow) error {
		return r.Scan(t)
	})
	return t, txb.Do(c)
}

type LatestPublicTagsResult struct {
	T        vocab.Type
	Received time.Time
}

func (d *DB) FetchLatestPublicTags(c util.Context, display []data.Tag, n int) ([]LatestPublicTagsResult, error) {
	// TODO
	return nil, nil
}

type RecentlyUpdatedThreadResult struct {
	First      vocab.Type
	MostRecent vocab.Type
}

func (d *DB) FetchMostRecentlyUpdatedThreads(c util.Context, t data.Tag, n, page int) ([]RecentlyUpdatedThreadResult, error) {
	// TODO
	return nil, nil
}

type ThreadMessages struct {
	Messages []vocab.Type
}

func (d *DB) FetchPaginatedMessagesInThread(c util.Context, id string, n, page int) (ThreadMessages, error) {
	// TODO
	return ThreadMessages{}, nil
}

func (d *DB) AddUserEmailValidationTask(c util.Context, userID, token string) error {
	txb := d.db.Begin()
	txb.ExecOneRow(d.pg.CreateUserSupplement(), userID, token, kUnvalidatedState)
	return txb.Do(c)
}

func (d *DB) MarkUserValidationEmailSent(c util.Context, userID string) error {
	txb := d.db.Begin()
	txb.ExecOneRow(d.pg.UpdateUserSupplement(), userID, kSentValidationChallengeState)
	return txb.Do(c)
}

func (d *DB) MarkUserVerified(c util.Context, token string) error {
	txb := d.db.Begin()
	txb.ExecOneRow(d.pg.UpdateUserVerified(), token, kValidatedState)
	return txb.Do(c)
}

func (d *DB) IsUserVerified(c util.Context, userID string) (bool, error) {
	// TODO: See if in-memory cache is needed for per-request latency
	var valid string
	txb := d.db.Begin()
	txb.QueryOneRow(d.pg.GetUserSupplement(), func(r app.SingleRow) error {
		return r.Scan(&valid)
	}, userID)
	return valid == kValidatedState, txb.Do(c)
}

func (d *DB) setApplicationState(c util.Context, k, v string) error {
	txb := d.db.Begin()
	txb.ExecOneRow(d.pg.SetApplicationStateKV(), k, v)
	return txb.Do(c)
}

func (d *DB) getApplicationState(c util.Context, k string) (string, error) {
	var value string
	txb := d.db.Begin()
	txb.QueryOneRow(d.pg.GetApplicationStateKV(), func(r app.SingleRow) error {
		return r.Scan(&value)
	}, k)
	return value, txb.Do(c)
}

func (d *DB) GetAuthoritativeCharacter(c util.Context) (string, error) {
	return d.getApplicationState(c, kAuthoritativeCharacterKey)
}

// TODO: Use
func (d *DB) SetAuthoritativeCharacter(c util.Context, v string) error {
	return d.setApplicationState(c, kAuthoritativeCharacterKey, v)
}

func (d *DB) GetCorporationManaged(c util.Context) (string, error) {
	return d.getApplicationState(c, kCorporationManagedKey)
}

// TODO: Use
func (d *DB) SetCorporationManaged(c util.Context, v string) error {
	return d.setApplicationState(c, kCorporationManagedKey, v)
}

func (d *DB) GetAlliance(c util.Context) (string, error) {
	return d.getApplicationState(c, kAllianceAssociationKey)
}

// TODO: Use
func (d *DB) SetAlliance(c util.Context, v string) error {
	return d.setApplicationState(c, kAllianceAssociationKey, v)
}

func (d *DB) GetExecutor(c util.Context) (string, error) {
	return d.getApplicationState(c, kExecutorCorporationKey)
}

// TODO: Use
func (d *DB) SetExecutor(c util.Context, v string) error {
	return d.setApplicationState(c, kExecutorCorporationKey, v)
}
