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

package features

import (
	"sort"

	"github.com/cjslep/dharma/internal/util"
	"github.com/cjslep/dharma/locales"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	CoreCorporationFeatureId = "core-corporation"
	CoreCalendarFeatureId    = "core-calendar"
	CoreMailFeatureId        = "core-mail"
	allFeatureIDs            = map[string]bool{
		CoreCorporationFeatureId: true,
		CoreCalendarFeatureId:    true,
		CoreMailFeatureId:        true,
	}
)

func allLocalizedFeatures(b *i18n.Bundle, langs ...string) ([]Feature, error) {
	m := locales.New(b, langs...)
	var err error
	return []Feature{
		{
			ID:          CoreCorporationFeatureId,
			Name:        util.MustPropagateString(m.FeatureCoreCorporationName, &err),
			Description: util.MustPropagateString(m.FeatureCoreCorporationDescription, &err),
			Scopes: []ScopeExplanation{
				{
					Scope:       "esi-corporations.read_corporation_membership.v1",
					Explanation: util.MustPropagateString(m.FeatureCoreCorporationReadCorporationMembershipScopeExplanation, &err),
				},
				{
					Scope:       "esi-corporations.read_fw_stats.v1",
					Explanation: util.MustPropagateString(m.FeatureCoreCorporationReadFactionWarfareStatsScopeExplanation, &err),
				},
				{
					Scope:       "esi-killmails.read_corporation_killmails.v1",
					Explanation: util.MustPropagateString(m.FeatureCoreCorporationReadCorporationKillmailsScopeExplanation, &err),
				},
				{
					Scope:       "esi-wallet.read_corporation_wallets.v1",
					Explanation: util.MustPropagateString(m.FeatureCoreCorporationReadCorporationWalletsScopeExplanation, &err),
				},
			},
			Required: true,
		},
		{
			ID:          CoreCalendarFeatureId,
			Name:        util.MustPropagateString(m.FeatureCoreCalendarName, &err),
			Description: util.MustPropagateString(m.FeatureCoreCalendarDescription, &err),
			Scopes: []ScopeExplanation{
				{
					Scope:       "esi-calendar.read_calendar_events.v1",
					Explanation: util.MustPropagateString(m.FeatureCoreReadCalendarEventsScopeExplanation, &err),
				},
				{
					Scope:       "esi-calendar.respond_calendar_events.v1",
					Explanation: util.MustPropagateString(m.FeatureCoreRespondCalendarEventsScopeExplanation, &err),
				},
			},
			Required: true,
		},
		{
			ID:          CoreMailFeatureId,
			Name:        util.MustPropagateString(m.FeatureCoreMailName, &err),
			Description: util.MustPropagateString(m.FeatureCoreMailDescription, &err),
			Scopes: []ScopeExplanation{
				{
					Scope:       "esi-mail.read_mail.v1",
					Explanation: util.MustPropagateString(m.FeatureCoreReadMailScopeExplanation, &err),
				},
				{
					Scope:       "esi-mail.send_mail.v1",
					Explanation: util.MustPropagateString(m.FeatureCoreSendMailScopeExplanation, &err),
				},
				{
					Scope:       "esi-mail.organize_mail.v1",
					Explanation: util.MustPropagateString(m.FeatureCoreOrganizeMailScopeExplanation, &err),
				},
			},
			Required: true,
		},
	}, err
}

type Feature struct {
	ID          string
	Name        string
	Description string
	Scopes      []ScopeExplanation
	Required    bool
}

type ScopeExplanation struct {
	Scope       Scope
	Explanation string
}

type List []Feature

func (l List) IDs() []string {
	s := make([]string, 0, len(l))
	for _, f := range l {
		s = append(s, f.ID)
	}
	// Deterministic list
	sort.Strings(s)
	return s
}

func (l List) Scopes() []Scope {
	// Deduplicate scopes
	m := make(map[Scope]bool, len(l))
	for i := range l {
		for _, se := range l[i].Scopes {
			m[se.Scope] = true
		}
	}
	// Put into slice
	s := make([]Scope, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	// Deterministic list
	sort.Sort(Scopes(s))
	return s
}

func (l List) ScopeExplanations() []ScopeExplanations {
	// Deduplicate scopes
	m := make(map[Scope][]string, len(l))
	for i := range l {
		for _, se := range l[i].Scopes {
			if x, ok := m[se.Scope]; ok {
				x = append(x, se.Explanation)
				m[se.Scope] = x
			} else {
				m[se.Scope] = []string{se.Explanation}
			}
		}
	}
	// Put into slice
	s := make([]ScopeExplanations, 0, len(m))
	for k, v := range m {
		s = append(s, ScopeExplanations{
			Scope:        k,
			Explanations: v,
		})
	}
	// Deterministic list
	sort.Sort(ScopeExplanationsList(s))
	return s
}
