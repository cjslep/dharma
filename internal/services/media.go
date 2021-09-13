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
	"database/sql"
	"image"
	"net/url"
	"time"

	"github.com/cjslep/dharma/esi"
	"github.com/cjslep/dharma/internal/db"
)

type SetEveMediaFn func(context.Context, int32, image.Image, time.Time) error
type GetEveMediaFn func(context.Context, int32) (image.Image, time.Time, error)

type PortraitMediaSize string

const (
	P64  PortraitMediaSize = "portrait64"
	P128 PortraitMediaSize = "portrait128"
	P256 PortraitMediaSize = "portrait256"
	P512 PortraitMediaSize = "portrait512"
)

func (p PortraitMediaSize) Apply(u *esi.PortraitURLs) *url.URL {
	switch p {
	default:
		fallthrough
	case P64:
		return u.Portrait64x64
	case P128:
		return u.Portrait128x128
	case P256:
		return u.Portrait256x256
	case P512:
		return u.Portrait512x512
	}
}

func (p PortraitMediaSize) GetEveMediaFn(d *db.DB) GetEveMediaFn {
	switch p {
	default:
		fallthrough
	case P64:
		return d.GetCharacterPortrait64
	case P128:
		return d.GetCharacterPortrait128
	case P256:
		return d.GetCharacterPortrait256
	case P512:
		return d.GetCharacterPortrait512
	}
}

func (p PortraitMediaSize) SetEveMediaFn(d *db.DB) SetEveMediaFn {
	switch p {
	default:
		fallthrough
	case P64:
		return d.SetCharacterPortrait64
	case P128:
		return d.SetCharacterPortrait128
	case P256:
		return d.SetCharacterPortrait256
	case P512:
		return d.SetCharacterPortrait512
	}
}

type CorpIconSize string

const (
	C64  CorpIconSize = "corp64"
	C128 CorpIconSize = "corp128"
	C256 CorpIconSize = "corp256"
)

func (c CorpIconSize) Apply(u *esi.CorpIconURLs) *url.URL {
	switch c {
	default:
		fallthrough
	case C64:
		return u.Icon64x64
	case C128:
		return u.Icon128x128
	case C256:
		return u.Icon256x256
	}
}

func (c CorpIconSize) GetEveMediaFn(d *db.DB) GetEveMediaFn {
	switch c {
	default:
		fallthrough
	case C64:
		return d.GetCorpIcon64
	case C128:
		return d.GetCorpIcon128
	case C256:
		return d.GetCorpIcon256
	}
}

func (c CorpIconSize) SetEveMediaFn(d *db.DB) SetEveMediaFn {
	switch c {
	default:
		fallthrough
	case C64:
		return d.SetCorpIcon64
	case C128:
		return d.SetCorpIcon128
	case C256:
		return d.SetCorpIcon256
	}
}

type AllianceIconSize string

const (
	A64  AllianceIconSize = "alliance64"
	A128 AllianceIconSize = "alliance128"
)

func (a AllianceIconSize) Apply(u *esi.AllianceIconURLs) *url.URL {
	switch a {
	default:
		fallthrough
	case A64:
		return u.Icon64x64
	case A128:
		return u.Icon128x128
	}
}

func (a AllianceIconSize) GetEveMediaFn(d *db.DB) GetEveMediaFn {
	switch a {
	default:
		fallthrough
	case A64:
		return d.GetAllianceIcon64
	case A128:
		return d.GetAllianceIcon128
	}
}

func (a AllianceIconSize) SetEveMediaFn(d *db.DB) SetEveMediaFn {
	switch a {
	default:
		fallthrough
	case A64:
		return d.SetAllianceIcon64
	case A128:
		return d.SetAllianceIcon128
	}
}

type Media struct {
	DB                    *db.DB
	ESIClient             *esi.Client
	DefaultExpiryDuration time.Duration
}

func (m *Media) GetPortrait(c context.Context, charID int32, size PortraitMediaSize) (image.Image, error) {
	getFn := size.GetEveMediaFn(m.DB)
	i, exp, err := getFn(c, charID)
	if err == sql.ErrNoRows || time.Now().After(exp) {
		// Refresh our cache of the image from Eve Online.
		purls, err := m.ESIClient.GetPortrait(c, charID)
		if err != nil {
			return nil, err
		}
		i, expires, err := esi.FetchEveOnlineImage(size.Apply(purls))
		if err != nil {
			return nil, err
		}
		if expires.IsZero() {
			expires = time.Now().Add(m.DefaultExpiryDuration)
		}
		// Update cached image
		setFn := size.SetEveMediaFn(m.DB)
		err = setFn(c, charID, i, expires)
		if err != nil {
			return nil, err
		}
		return i, nil
	} else if err == nil {
		// Our cache version is OK.
		return i, nil
	} else {
		return nil, err
	}
}

func (m *Media) GetCorporationIcon(c context.Context, corpID int32, size CorpIconSize) (image.Image, error) {
	getFn := size.GetEveMediaFn(m.DB)
	i, exp, err := getFn(c, corpID)
	if err == sql.ErrNoRows || time.Now().After(exp) {
		// Refresh our cache of the image from Eve Online.
		curls, err := m.ESIClient.GetCorporationIcon(c, corpID)
		if err != nil {
			return nil, err
		}
		i, expires, err := esi.FetchEveOnlineImage(size.Apply(curls))
		if err != nil {
			return nil, err
		}
		if expires.IsZero() {
			expires = time.Now().Add(m.DefaultExpiryDuration)
		}
		// Update cached image
		setFn := size.SetEveMediaFn(m.DB)
		err = setFn(c, corpID, i, expires)
		if err != nil {
			return nil, err
		}
		return i, nil
	} else if err == nil {
		// Our cache version is OK.
		return i, nil
	} else {
		return nil, err
	}
}

func (m *Media) GetAllianceIcon(c context.Context, aID int32, size AllianceIconSize) (image.Image, error) {
	getFn := size.GetEveMediaFn(m.DB)
	i, exp, err := getFn(c, aID)
	if err == sql.ErrNoRows || time.Now().After(exp) {
		// Refresh our cache of the image from Eve Online.
		curls, err := m.ESIClient.GetAllianceIcon(c, aID)
		if err != nil {
			return nil, err
		}
		i, expires, err := esi.FetchEveOnlineImage(size.Apply(curls))
		if err != nil {
			return nil, err
		}
		if expires.IsZero() {
			expires = time.Now().Add(m.DefaultExpiryDuration)
		}
		// Update cached image
		setFn := size.SetEveMediaFn(m.DB)
		err = setFn(c, aID, i, expires)
		if err != nil {
			return nil, err
		}
		return i, nil
	} else if err == nil {
		// Our cache version is OK.
		return i, nil
	} else {
		return nil, err
	}
}

func (m *Media) GetImage(id string) {
	// TODO
}

func (m *Media) SaveImage(id string, img image.Image) {
	// TODO
}

func (m *Media) DeleteImage(id string) {
	// TODO
}
