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
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/cjslep/dharma/esi"
	"github.com/cjslep/dharma/internal/data"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/go-fed/apcore/app"
	"github.com/pkg/errors"
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
	pg    postgres
	cache *cache
}

func New(db app.Database, schema string) *DB {
	return &DB{
		db:    db,
		pg:    newPostgres(schema),
		cache: newCache(),
	}
}

func (d *DB) SetEvePublicKeys(c context.Context, o *esi.OAuthKeysMetadata) error {
	txb := d.db.Begin()
	txb.Exec(d.pg.AddEvePublicKey(), o)
	return txb.Do(c)
}

func (d *DB) GetEvePublicKeys(c context.Context) (*esi.OAuthKeysMetadata, error) {
	// TODO: See if in-memory cache is needed for per-request latency
	o := &esi.OAuthKeysMetadata{}
	txb := d.db.Begin()
	txb.QueryOneRow(d.pg.GetLatestEvePublicKey(), func(r app.SingleRow) error {
		return r.Scan(o)
	})
	return o, txb.Do(c)
}

const (
	tokenOKState      = "ok"
	tokenRescopeState = "rescope"
)

func (d *DB) SetEveTokens(c context.Context, userID string, t *esi.Tokens) error {
	txb := d.db.Begin()
	txb.ExecOneRow(d.pg.SetEveToken(), userID, t.CID, t, tokenOKState)
	return txb.Do(c)
}

func (d *DB) GetEveToken(c context.Context, charID int32) (*esi.Tokens, error) {
	// TODO: Be more intelligent about refreshing keys and detecting cacheability
	t := &esi.Tokens{}
	txb := d.db.Begin()
	txb.QueryOneRow(d.pg.GetEveToken(), func(r app.SingleRow) error {
		return r.Scan(t)
	}, charID)
	return t, txb.Do(c)
}

func (d *DB) MarkAllTokensNeedRescope(c context.Context) error {
	txb := d.db.Begin()
	txb.Exec(d.pg.MarkAllTokensWithState(), tokenRescopeState)
	return txb.Do(c)
}

type UserToken struct {
	UserID string
	T      *esi.Tokens
}

func (d *DB) GetExpiringEveTokensWithin(c context.Context, period time.Duration) ([]*UserToken, error) {
	nHours := int(math.Ceil(period.Hours()))
	var ut []*UserToken
	txb := d.db.Begin()
	txb.Query(d.pg.GetExpiringEveTokensWithin(), func(r app.SingleRow) error {
		var u UserToken
		err := r.Scan(&(u.UserID), u.T)
		if err != nil {
			return err
		}
		ut = append(ut, &u)
		return nil
	}, fmt.Sprintf("%d hours", nHours))
	return ut, txb.Do(c)
}

func (d *DB) GetEveCharactersForUser(c context.Context, userID string) ([]int32, error) {
	var ids []int32
	txb := d.db.Begin()
	txb.Query(d.pg.GetEveCharactersForUser(), func(r app.SingleRow) error {
		var i int32
		err := r.Scan(&i)
		if err != nil {
			return err
		}
		ids = append(ids, i)
		return nil
	}, userID)
	return ids, txb.Do(c)
}

func (d *DB) HasCharacterForUser(c context.Context, userID string, charID int32) (bool, error) {
	txb := d.db.Begin()
	var ok bool
	txb.QueryOneRow(d.pg.HasCharacterForUser(), func(r app.SingleRow) error {
		err := r.Scan(&ok)
		return err
	}, userID, charID)
	return ok, txb.Do(c)
}

type LatestPublicTagsResult struct {
	T        vocab.Type
	Received time.Time
}

func (d *DB) FetchLatestPublicTags(c context.Context, display []data.Tag, n int) ([]LatestPublicTagsResult, error) {
	// TODO
	return nil, nil
}

type RecentlyUpdatedThreadResult struct {
	First      vocab.Type
	MostRecent vocab.Type
}

func (d *DB) FetchMostRecentlyUpdatedThreads(c context.Context, t data.Tag, n, page int) ([]RecentlyUpdatedThreadResult, error) {
	// TODO
	return nil, nil
}

type ThreadMessages struct {
	Messages []vocab.Type
}

func (d *DB) FetchPaginatedMessagesInThread(c context.Context, id string, n, page int) (ThreadMessages, error) {
	// TODO
	return ThreadMessages{}, nil
}

func (d *DB) AddUserEmailValidationTask(c context.Context, userID, token string) error {
	txb := d.db.Begin()
	txb.ExecOneRow(d.pg.CreateUserSupplement(), userID, token, kUnvalidatedState)
	return txb.Do(c)
}

func (d *DB) MarkUserValidationEmailSent(c context.Context, userID string) error {
	txb := d.db.Begin()
	txb.ExecOneRow(d.pg.UpdateUserSupplement(), userID, kSentValidationChallengeState)
	return txb.Do(c)
}

func (d *DB) MarkUserVerified(c context.Context, token string) error {
	txb := d.db.Begin()
	txb.ExecOneRow(d.pg.UpdateUserVerified(), token, kValidatedState)
	return txb.Do(c)
}

func (d *DB) IsUserVerified(c context.Context, userID string) (bool, error) {
	// TODO: See if in-memory cache is needed for per-request latency
	var valid string
	txb := d.db.Begin()
	txb.QueryOneRow(d.pg.GetUserSupplement(), func(r app.SingleRow) error {
		return r.Scan(&valid)
	}, userID)
	return valid == kValidatedState, txb.Do(c)
}

func (d *DB) setApplicationState(c context.Context, k, v string) error {
	txb := d.db.Begin()
	txb.ExecOneRow(d.pg.SetApplicationStateKV(), k, v)
	return txb.Do(c)
}

func (d *DB) setApplicationStateTx(tx app.TxBuilder, k, v string) {
	tx.ExecOneRow(d.pg.SetApplicationStateKV(), k, v)
}

func (d *DB) getApplicationState(c context.Context, k string) (string, error) {
	var value string
	txb := d.db.Begin()
	txb.QueryOneRow(d.pg.GetApplicationStateKV(), func(r app.SingleRow) error {
		return r.Scan(&value)
	}, k)
	return value, txb.Do(c)
}

func (d *DB) getApplicationStateAsInt32(c context.Context, k string) (int32, error) {
	s, err := d.getApplicationState(c, k)
	if err != nil {
		return 0, err
	}
	return stateValueToInt32(s)
}

func (d *DB) setApplicationStateAsInt32(c context.Context, k string, v int32) error {
	return d.setApplicationState(c, k, int32ToStateValue(v))
}

func (d *DB) GetAuthoritativeCharacter(c context.Context) (int32, error) {
	return d.getApplicationStateAsInt32(c, kAuthoritativeCharacterKey)
}

func (d *DB) SetAuthoritativeCharacter(c context.Context, v int32) error {
	return d.setApplicationStateAsInt32(c, kAuthoritativeCharacterKey, v)
}

func (d *DB) GetCorporationManaged(c context.Context) (int32, error) {
	return d.getApplicationStateAsInt32(c, kCorporationManagedKey)
}

func (d *DB) SetCorporationManaged(c context.Context, v int32) error {
	return d.setApplicationStateAsInt32(c, kCorporationManagedKey, v)
}

func (d *DB) GetAlliance(c context.Context) (int32, error) {
	return d.getApplicationStateAsInt32(c, kAllianceAssociationKey)
}

func (d *DB) SetAlliance(c context.Context, v int32) error {
	return d.setApplicationStateAsInt32(c, kAllianceAssociationKey, v)
}

func (d *DB) GetExecutor(c context.Context) (int32, error) {
	return d.getApplicationStateAsInt32(c, kExecutorCorporationKey)
}

func (d *DB) SetExecutor(c context.Context, v int32) error {
	return d.setApplicationStateAsInt32(c, kExecutorCorporationKey, v)
}

func (d *DB) setApplicationStateAsBoolTx(tx app.TxBuilder, k string, v bool) {
	d.setApplicationStateTx(tx, k, boolToStateValue(v))
}

const (
	featurePrefix = "feature-"
)

func (d *DB) addFeaturePrefix(s string) string {
	return featurePrefix + s
}

func (d *DB) removeFeaturePrefix(s string) string {
	return strings.TrimPrefix(s, featurePrefix)
}

func (d *DB) setFeatureEnabledTx(tx app.TxBuilder, featureID string) {
	d.setApplicationStateAsBoolTx(tx, d.addFeaturePrefix(featureID), true)
}

func (d *DB) setFeatureDisabledTx(tx app.TxBuilder, featureID string) {
	d.setApplicationStateAsBoolTx(tx, d.addFeaturePrefix(featureID), false)
}

func (d *DB) SetFeaturesEnabledDisabled(c context.Context, enableIDs, disableIDs []string) error {
	tx := d.db.Begin()
	for _, id := range enableIDs {
		d.setFeatureEnabledTx(tx, id)
	}
	for _, id := range disableIDs {
		d.setFeatureDisabledTx(tx, id)
	}
	return tx.Do(c)
}

func (d *DB) GetEnabledFeatureIDs(c context.Context) (ids []string, err error) {
	txb := d.db.Begin()
	txb.Query(d.pg.GetApplicationStateValueAndPrefix(), func(r app.SingleRow) error {
		var id string
		if err := r.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
		return nil
	}, featurePrefix, boolToStateValue(true))
	err = txb.Do(c)
	return
}

func int32ToStateValue(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func stateValueToInt32(s string) (int32, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	return int32(i), err
}

func boolToStateValue(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

func stateValueToBool(s string) (bool, error) {
	if s == "true" {
		return true, nil
	} else if s == "false" {
		return false, nil
	} else {
		return false, errors.Errorf("error parsing bool for state value: %s", s)
	}
}
