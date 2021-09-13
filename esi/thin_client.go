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

package esi

import (
	"context"
	"net/http"
	"time"

	"github.com/cjslep/dharma/esi/client"
	"github.com/cjslep/dharma/esi/client/alliance"
	"github.com/cjslep/dharma/esi/client/character"
	"github.com/cjslep/dharma/esi/client/corporation"
	"github.com/cjslep/dharma/esi/client/search"
	"golang.org/x/text/language"
)

var (
	server = "tranquility"
)

type ThinClient struct {
	ESIClient *client.ESI
	Timeout   time.Duration
	Client    *http.Client
}

// search is a thin wrapper for ESI Search.
func (e *ThinClient) search(c context.Context, q string, cats []string, isStrict bool, l language.Tag) (*search.GetSearchOKBody, error) {
	p := search.NewGetSearchParams()
	lang := l.String()
	p.WithTimeout(e.Timeout).
		WithContext(c).
		WithHTTPClient(e.Client).
		WithCategories(cats).
		WithDatasource(&server).
		WithLanguage(&lang).
		WithSearch(q).
		WithStrict(&isStrict)
	resp, err := e.ESIClient.Search.GetSearch(p)
	if err != nil {
		return nil, err
	}
	return resp.GetPayload(), nil
}

// corporation is a thin wrapper for ESI corporation.
func (e *ThinClient) corporation(c context.Context, id int32) (*corporation.GetCorporationsCorporationIDOKBody, error) {
	p := corporation.NewGetCorporationsCorporationIDParams()
	p.WithTimeout(e.Timeout).
		WithContext(c).
		WithHTTPClient(e.Client).
		WithCorporationID(id)
	resp, err := e.ESIClient.Corporation.GetCorporationsCorporationID(p)
	if err != nil {
		return nil, err
	}
	return resp.GetPayload(), nil
}

// alliance is a thin wrapper for ESI alliance.
func (e *ThinClient) alliance(c context.Context, id int32) (*alliance.GetAlliancesAllianceIDOKBody, error) {
	p := alliance.NewGetAlliancesAllianceIDParams()
	p.WithTimeout(e.Timeout).
		WithContext(c).
		WithHTTPClient(e.Client).
		WithAllianceID(id)
	resp, err := e.ESIClient.Alliance.GetAlliancesAllianceID(p)
	if err != nil {
		return nil, err
	}
	return resp.GetPayload(), nil
}

// character is a thin wrapper for ESI character.
func (e *ThinClient) character(c context.Context, id int32) (*character.GetCharactersCharacterIDOKBody, error) {
	p := character.NewGetCharactersCharacterIDParams()
	p.WithTimeout(e.Timeout).
		WithContext(c).
		WithHTTPClient(e.Client).
		WithCharacterID(id)
	resp, err := e.ESIClient.Character.GetCharactersCharacterID(p)
	if err != nil {
		return nil, err
	}
	return resp.GetPayload(), nil
}

// characterPortrait is a thin wrapper for ESI character portrait.
func (e *ThinClient) characterPortrait(c context.Context, id int32) (*character.GetCharactersCharacterIDPortraitOKBody, error) {
	p := character.NewGetCharactersCharacterIDPortraitParams()
	p.WithTimeout(e.Timeout).
		WithContext(c).
		WithHTTPClient(e.Client).
		WithCharacterID(id)
	resp, err := e.ESIClient.Character.GetCharactersCharacterIDPortrait(p)
	if err != nil {
		return nil, err
	}
	return resp.GetPayload(), nil
}

// corporationIcon is a thin wrapper for ESI corporation icon.
func (e *ThinClient) corporationIcon(c context.Context, id int32) (*corporation.GetCorporationsCorporationIDIconsOKBody, error) {
	p := corporation.NewGetCorporationsCorporationIDIconsParams()
	p.WithTimeout(e.Timeout).
		WithContext(c).
		WithHTTPClient(e.Client).
		WithCorporationID(id)
	resp, err := e.ESIClient.Corporation.GetCorporationsCorporationIDIcons(p)
	if err != nil {
		return nil, err
	}
	return resp.GetPayload(), nil
}

// allianceIcon is a thin wrapper for ESI alliance icon.
func (e *ThinClient) allianceIcon(c context.Context, id int32) (*alliance.GetAlliancesAllianceIDIconsOKBody, error) {
	p := alliance.NewGetAlliancesAllianceIDIconsParams()
	p.WithTimeout(e.Timeout).
		WithContext(c).
		WithHTTPClient(e.Client).
		WithAllianceID(id)
	resp, err := e.ESIClient.Alliance.GetAlliancesAllianceIDIcons(p)
	if err != nil {
		return nil, err
	}
	return resp.GetPayload(), nil
}
