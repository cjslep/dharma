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
	"github.com/cjslep/dharma/internal/db"
	"github.com/go-fed/apcore/util"
	"github.com/pkg/errors"
)

// State solely manages the state of the dharma application in the face of
// complexity.
//
// The complexity arises because the problem space is:
// 1) There are X accounts on dharma, a subset of which are admin,
// 2) There are numerous EVE Online character(s) associated with accounts on
//    dharma,
// 3) There is exactly one CEO (EVE Online character) of the corporation using
//    dharma,
// 4) Optionally, there is an executor corporation of an alliance, where an
//    alliance is a collection of one or more alliances
// 5) Finally, while not a game mechanic, there is the concept of coalitions
//    which already are effecively alliances that federate for a common cause
//
// In the spirit of decentralized leadership, it simplifies the problem space
// by enforcing the following constraints to simplify software administration
// and push power downwards:
//
// 0) A dharma account with no characters or only non-CEO characters does not
//    enable management of any corporation.
// 1) Dharma manages ONE and ONLY one corporation.
// 2) A dharma administrative account, paired with an EVE Online character that
//    is a CEO, makes that CEO's corporation "eligible" to be managed by dharma
// 3) Corporations managed by dharma, that are also in alliances, have automatic
//    federating behaviors that occur but only towards peer software managing
//    corporations that are also in the same alliance.
// 4) Dharma does nothing special for coalitions, these are covered by the
//    existing ad-hoc federation controls and states.
//
// Note that #0 and #2 imply that there can be many administrative accounts in
// dharma, but only one of which requires a CEO EVE Online character. This
// ensures that the in-game CEO always has the maximum tools necessary to
// administrate dharma, yet retains the option to delegate devops to others.
type State struct {
	state appState
	db    *db.DB
}

type appState string

const (
	// unmanagedState is the initial state where there is no corporation
	// being managed by dharma.
	unmanagedState appState = "unmanaged"
	// managedIndyCorpState is a state where there is a corporation being
	// managed by dharma, and it is not a part of an alliance.
	managedIndyCorpState = "managed_independent_corp"
	// managedAllianceCorpState is a state where there is a non-executor
	// corp being managed by dharma that is also a part of an alliance.
	managedAllianceCorpState = "managed_alliance_corp"
	// managedExecutorCorpState is a state where there is an executor
	// corp being managed by dharma that is leading an alliance.
	managedExecutorCorpState = "managed_executor_corp"
)

func NewState(c util.Context, db *db.DB) (*State, error) {
	s := &State{
		db: db,
	}
	if h, err := db.GetAuthoritativeCharacter(c); err != nil || len(h) == 0 {
		s.state = unmanagedState
	} else if r, err := db.GetCorporationManaged(c); err != nil || len(r) == 0 {
		s.state = unmanagedState
	} else if a, err := db.GetAlliance(c); err != nil || len(a) == 0 {
		s.state = managedIndyCorpState
	} else if x, err := db.GetExecutor(c); err != nil || len(x) == 0 || x != r {
		s.state = managedAllianceCorpState
	} else if x == r {
		s.state = managedExecutorCorpState
	} else {
		return nil, errors.Errorf("cannot initialize application state: char=%s, corp=%s, alli=%s, exec=%s", h, r, a, x)
	}
	return s, nil
}

// RequiresCorpToBeManaged determines if the current state requires an admin
// account to log in with a CEO character to manage that corporation.
func (s *State) RequiresCorpToBeManaged() bool {
	return s.state == unmanagedState
}

// ShouldCorpReceiveAllianceData determines if the current state requires the
// corporation to receive alliance-level data from peers automatically.
func (s *State) ShouldCorpReceiveAllianceData() bool {
	switch s.state {
	case managedAllianceCorpState:
		fallthrough
	case managedExecutorCorpState:
		return true
	default:
		return false
	}
}

// ShouldCorpSendAllianceData determines if the current state requires the
// corporation to send alliance-level data to peers automatically.
func (s *State) ShouldCorpSendAllianceData() bool {
	switch s.state {
	case managedAllianceCorpState:
		fallthrough
	case managedExecutorCorpState:
		return true
	default:
		return false
	}
}

func (s *State) ChooseCorporation(c util.Context, userID string, corpID int32) error {
	// TODO: Check that userID has an ESI character token associated with them
	// TODO: Check that the character ID of one of the ESI characters is the corp CEO
	var charStr, corpStr, allianceStr, execStr string
	if err := s.db.SetAuthoritativeCharacter(c, charStr); err != nil {
		return err
	}
	if err := s.db.SetCorporationManaged(c, corpStr); err != nil {
		return err
	}
	if err := s.db.SetAlliance(c, allianceStr); err != nil {
		return err
	}
	if err := s.db.SetExecutor(c, execStr); err != nil {
		return err
	}
	return nil
}
