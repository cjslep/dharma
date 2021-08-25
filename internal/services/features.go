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

	"github.com/cjslep/dharma/internal/db"
	"github.com/cjslep/dharma/internal/features"
	"golang.org/x/text/language"
)

type Features struct {
	DB *db.DB
	E  *features.Engine
}

func (f *Features) DiffChangeEnabled(ctx context.Context, enableIDs, disableIDs []string) (added, removed []features.Scope, err error) {
	// 0. Validate feature IDs
	if err = f.E.ValidateFeatureIDs(enableIDs); err != nil {
		return
	}
	if err = f.E.ValidateFeatureIDs(disableIDs); err != nil {
		return
	}
	// 1. Diff against existing state to determine scope change(s)
	var ids []string
	ids, err = f.DB.GetEnabledFeatureIDs(ctx)
	if err != nil {
		return
	}
	added, removed, err = f.E.DiffScopes(ids, enableIDs, disableIDs, language.English)
	return
}

func (f *Features) ChangeEnabled(ctx context.Context, enableIDs, disableIDs []string) error {
	// 0. Validate feature IDs
	if err := f.E.ValidateFeatureIDs(enableIDs); err != nil {
		return err
	}
	if err := f.E.ValidateFeatureIDs(disableIDs); err != nil {
		return err
	}
	// 1. Diff against existing state to determine scope change(s)
	ids, err := f.DB.GetEnabledFeatureIDs(ctx)
	if err != nil {
		return err
	}
	added, removed, err := f.E.DiffScopes(ids, enableIDs, disableIDs, language.English)
	if err != nil {
		return err
	}
	scopeChange := len(added) > 0 || len(removed) > 0
	// 2. Change the features
	if err := f.DB.SetFeaturesEnabledDisabled(ctx, enableIDs, disableIDs); err != nil {
		return err
	}
	// 3. If there is a scope change, mark tokens as needing rescope
	if scopeChange {
		if err := f.DB.MarkAllTokensNeedRescope(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (f *Features) GetAdminCEOInitialFeatures(ctx context.Context, langs ...language.Tag) ([]features.Feature, error) {
	return f.E.GetRequiredFeatures(langs...)
}

func (f *Features) GetEnabledFeatures(ctx context.Context, langs ...language.Tag) ([]features.Feature, error) {
	ids, err := f.DB.GetEnabledFeatureIDs(ctx)
	if err != nil {
		return nil, err
	}
	return f.E.GetFeatures(ids, langs...)
}

func (f *Features) GetByIDs(ctx context.Context, ids []string, langs ...language.Tag) ([]features.Feature, error) {
	return f.E.GetFeatures(ids, langs...)
}

func (f *Features) GetAllFeatures(ctx context.Context, langs ...language.Tag) ([]features.Feature, error) {
	return f.E.GetAllFeatures(langs...)
}
