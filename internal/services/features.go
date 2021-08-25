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

func (f *Features) EnableFeature(ctx context.Context, id string) error {
	if err := f.E.ValidateFeatureID(id); err != nil {
		return err
	}
	// TODO: Invalidate existing tokens if scope change required
	return f.DB.SetFeatureEnabled(ctx, id)
}

func (f *Features) DisableFeature(ctx context.Context, id string) error {
	if err := f.E.ValidateFeatureID(id); err != nil {
		return err
	}
	// TODO: Invalidate existing tokens if scope change required
	return f.DB.SetFeatureDisabled(ctx, id)
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
