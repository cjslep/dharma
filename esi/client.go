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
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type Client struct {
	t *ThinClient
}

func New(t *ThinClient) *Client {
	return &Client{
		t: t,
	}
}

type Alliance struct {
	ID       int32 `json:"id,omitempty"`
	Hydrated bool  `json:"-"`
	// Only set if hydrated
	Name         string       `json:"name,omitempty"`
	Ticker       string       `json:"ticker,omitempty"`
	CreationCorp *Corporation `json:"creation_corporation,omitempty"`
	Creator      *Character   `json:"creator,omitempty"`
	Founded      time.Time    `json:"founded,omitempty"`
	Executor     *Corporation `json:"executor_corproration,omitempty"`
	Faction      *Faction     `json:"faction,omitempty"`
}

type Corporation struct {
	ID       int32 `json:"id,omitempty"`
	Hydrated bool  `json:"-"`
	// Only set if hydrated
	Name        string     `json:"name,omitempty"`
	Ticker      string     `json:"ticker,omitempty"`
	CEO         *Character `json:"ceo,omitempty"`
	Creator     *Character `json:"creator,omitempty"`
	Alliance    *Alliance  `json:"alliance,omitempty"`
	Founded     time.Time  `json:"founded,omitempty"`
	Description string     `json:"description,omitempty"`
	Faction     *Faction   `json:"faction,omitempty"`
	Home        *Station   `json:"home_station,omitempty"`
	NMembers    int32      `json:"number_members,omitempty"`
	Shares      int64      `json:"shares,omitempty"`
	TaxRate     float64    `json:"tax_rate,omitempty"`
	URL         *url.URL   `json:"url,omitempty"`
	WarEligible bool       `json:"war_eligible,omitempty"`
}

type Character struct {
	ID       int32 `json:"id,omitempty"`
	Hydrated bool  `json:"-"`
	// Only set if hydrated
	Name           string       `json:"name,omitempty"`
	Title          string       `json:"title,omitempty"`
	Corporation    *Corporation `json:"corporation,omitempty"`
	Alliance       *Alliance    `json:"alliance,omitempty"`
	Birthday       time.Time    `json:"birthday,omitempty"`
	Description    string       `json:"description,omitempty"`
	SecurityStatus float64      `json:"security_status,omitempty"`
	Faction        *Faction     `json:"faction,omitempty"`
	Gender         string       `json:"gender,omitempty"`
	Ancestry       *Ancestry    `json:"ancestry,omitempty"`
	Bloodline      *Bloodline   `json:"bloodline,omitempty"`
	Race           *Race        `json:"race,omitempty"`
}

type Faction struct {
	ID int32 `json:"id,omitempty"`
}

type Station struct {
	ID int32 `json:"id,omitempty"`
}

type Ancestry struct {
	ID int32 `json:"id,omitempty"`
}

type Bloodline struct {
	ID int32 `json:"id,omitempty"`
}

type Race struct {
	ID int32 `json:"id,omitempty"`
}

// SearchCorp returns a slice of Corporations matching the search query
//
// Hydrates the Corporation, any associated alliance, the CEO, and Creator.
func (x *Client) SearchCorp(ctx context.Context, q string, l language.Tag) ([]*Corporation, error) {
	p, err := x.t.search(ctx, q, []string{"corporation"}, false, l)
	if err != nil {
		return nil, err
	}
	if len(p.Corporation) == 0 {
		return nil, nil
	}
	var wg sync.WaitGroup
	wg.Add(len(p.Corporation))
	corps := make([]*Corporation, len(p.Corporation))
	errs := make([]error, len(p.Corporation))
	for idx, c := range p.Corporation {
		go func(idx int, corpID int32) {
			defer wg.Done()
			corp, err := x.Corporation(ctx, corpID)
			if err != nil {
				errs[idx] = errors.Wrapf(err, "error fetching corporation id: %d", corpID)
				return
			}
			corps[idx] = corp
		}(idx, c)
	}
	wg.Wait()
	err = toErrors(errs)
	if err != nil {
		return nil, err
	}
	return corps, nil
}

// Character obtains information about the character.
//
// Hydrates the Corporation, and any associated alliance.
func (x *Client) Character(ctx context.Context, id int32) (*Character, error) {
	c, err := x.hydrateCharacter(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.Corporation != nil {
		cID := c.Corporation.ID
		c.Corporation, err = x.hydrateCorporation(ctx, cID)
		if err != nil {
			return nil, errors.Wrapf(err, "error fetching corporation id: %d", cID)
		}
	}
	if c.Alliance != nil {
		aID := c.Alliance.ID
		c.Alliance, err = x.hydrateAlliance(ctx, aID)
		if err != nil {
			return nil, errors.Wrapf(err, "error fetching alliance id: %d", aID)
		}
	}
	return c, nil
}

// Characters is the same as Character but in parallel for multiple.
func (x *Client) Characters(ctx context.Context, ids []int32) ([]*Character, error) {
	var wg sync.WaitGroup
	wg.Add(len(ids))
	chars := make([]*Character, len(ids))
	errs := make([]error, len(ids))
	for idx, id := range ids {
		go func(idx int, id int32) {
			defer wg.Done()
			ch, err := x.Character(ctx, id)
			if err != nil {
				errs[idx] = err
			} else {
				chars[idx] = ch
			}
		}(idx, id)
	}
	wg.Wait()
	err := toErrors(errs)
	if err != nil {
		return nil, err
	}
	return chars, nil
}

// Corporation obtains information about the Corporation.
//
// Hydrates the CEO, Creator, and any associated alliance.
func (x *Client) Corporation(ctx context.Context, id int32) (*Corporation, error) {
	corp, err := x.hydrateCorporation(ctx, id)
	if err != nil {
		return nil, err
	}
	if corp.Alliance != nil {
		aID := corp.Alliance.ID
		corp.Alliance, err = x.hydrateAlliance(ctx, aID)
		if err != nil {
			return nil, errors.Wrapf(err, "error fetching alliance id: %d", aID)
		}
	}
	if corp.CEO != nil {
		cID := corp.CEO.ID
		corp.CEO, err = x.hydrateCharacter(ctx, cID)
		if err != nil {
			return nil, errors.Wrapf(err, "error fetching ceo id: %d", cID)
		}
	}
	if corp.Creator != nil {
		cID := corp.Creator.ID
		corp.Creator, err = x.hydrateCharacter(ctx, cID)
		if err != nil {
			return nil, errors.Wrapf(err, "error fetching creator id: %d", cID)
		}
	}
	return corp, nil
}

func (x *Client) hydrateCorporation(ctx context.Context, id int32) (*Corporation, error) {
	p, err := x.t.corporation(ctx, id)
	if err != nil {
		return nil, err
	}
	var c Corporation
	c.ID = id
	c.Hydrated = true
	if p.Name != nil {
		c.Name = *p.Name
	}
	if p.Ticker != nil {
		c.Ticker = *p.Ticker
	}
	if p.CeoID != nil && *p.CeoID != 0 {
		c.CEO = &Character{
			ID: *p.CeoID,
		}
	}
	if p.CreatorID != nil && *p.CreatorID != 0 {
		c.Creator = &Character{
			ID: *p.CreatorID,
		}
	}
	if p.AllianceID != 0 {
		c.Alliance = &Alliance{
			ID: p.AllianceID,
		}
	}
	c.Founded = time.Time(p.DateFounded)
	c.Description = p.Description
	if p.FactionID != 0 {
		c.Faction = &Faction{
			ID: p.FactionID,
		}
	}
	if p.HomeStationID != 0 {
		c.Home = &Station{
			ID: p.HomeStationID,
		}
	}
	if p.MemberCount != nil {
		c.NMembers = *p.MemberCount
	}
	c.Shares = p.Shares
	if p.TaxRate != nil {
		c.TaxRate = float64(*p.TaxRate)
	}
	if p.URL != "" &&
		p.URL != "http://" &&
		p.URL != "https://" {
		c.URL, err = url.Parse(p.URL)
		if err != nil {
			return nil, err
		}
	}
	c.WarEligible = p.WarEligible
	return &c, nil
}

func (x *Client) hydrateAlliance(ctx context.Context, id int32) (*Alliance, error) {
	p, err := x.t.alliance(ctx, id)
	if err != nil {
		return nil, err
	}
	var a Alliance
	a.ID = id
	a.Hydrated = true
	if p.Name != nil {
		a.Name = *p.Name
	}
	if p.Ticker != nil {
		a.Ticker = *p.Ticker
	}
	if p.CreatorCorporationID != nil && *p.CreatorCorporationID != 0 {
		a.CreationCorp = &Corporation{
			ID: *p.CreatorCorporationID,
		}
	}
	if p.CreatorID != nil && *p.CreatorID != 0 {
		a.Creator = &Character{
			ID: *p.CreatorID,
		}
	}
	if p.DateFounded != nil {
		a.Founded = time.Time(*p.DateFounded)
	}
	if p.ExecutorCorporationID != 0 {
		a.Executor = &Corporation{
			ID: p.ExecutorCorporationID,
		}
	}
	if p.FactionID != 0 {
		a.Faction = &Faction{
			ID: p.FactionID,
		}
	}
	return &a, nil
}

func (x *Client) hydrateCharacter(ctx context.Context, id int32) (*Character, error) {
	p, err := x.t.character(ctx, id)
	if err != nil {
		return nil, err
	}
	var c Character
	c.ID = id
	c.Hydrated = true
	// Only set if hydrated
	if p.Name != nil {
		c.Name = *p.Name
	}
	c.Title = p.Title
	if p.CorporationID != nil && *p.CorporationID != 0 {
		c.Corporation = &Corporation{
			ID: *p.CorporationID,
		}
	}
	if p.AllianceID != 0 {
		c.Alliance = &Alliance{
			ID: p.AllianceID,
		}
	}
	if p.Birthday != nil {
		c.Birthday = time.Time(*p.Birthday)
	}
	c.Description = p.Description
	if p.SecurityStatus != nil {
		c.SecurityStatus = float64(*p.SecurityStatus)
	}
	if p.FactionID != 0 {
		c.Faction = &Faction{
			ID: p.FactionID,
		}
	}
	if p.Gender != nil && *p.Gender != "" {
		c.Gender = *p.Gender
	}
	if p.AncestryID != 0 {
		c.Ancestry = &Ancestry{
			ID: p.AncestryID,
		}
	}
	if p.BloodlineID != nil && *p.BloodlineID != 0 {
		c.Bloodline = &Bloodline{
			ID: *p.BloodlineID,
		}
	}
	if p.RaceID != nil && *p.RaceID != 0 {
		c.Race = &Race{
			ID: *p.RaceID,
		}
	}
	return &c, nil
}

var _ error = &manyErrors{}

type manyErrors []error

func toErrors(nilable []error) error {
	notNilErrs := make([]error, 0, len(nilable))
	for _, err := range nilable {
		if err != nil {
			notNilErrs = append(notNilErrs, err)
		}
	}
	if len(notNilErrs) == 0 {
		return nil
	}
	e := manyErrors(notNilErrs)
	return &e
}

func (e *manyErrors) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d errors occurred:", len(*e))
	for _, err := range *e {
		fmt.Fprintf(&b, "\n> %s", err)
	}
	return b.String()
}
