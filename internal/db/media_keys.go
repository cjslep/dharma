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

const (
	// Kinds of media stored in the database for the Eve Media Table.
	//
	// Be explicit about the values used and avoid using iota syntactic
	// sugar, to ease maintenance as kinds are introduced or removed.
	mediaKindCharacterPortrait64  = 1
	mediaKindCharacterPortrait128 = 2
	mediaKindCharacterPortrait256 = 3
	mediaKindCharacterPortrait512 = 4
	mediaKindCorpIcon64           = 5
	mediaKindCorpIcon128          = 6
	mediaKindCorpIcon256          = 7
	mediaKindAllianceIcon64       = 8
	mediaKindAllianceIcon128      = 9
)
