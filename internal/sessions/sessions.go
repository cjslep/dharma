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

package sessions

import (
	"github.com/go-fed/apcore/app"
)

const (
	esiOAuth2State          = "dharma-esi-oauth2-state"
	dharmaCharacterSelected = "dharma-character-selected"
)

func SetESIOAuth2State(k app.Session, state string) {
	k.Set(esiOAuth2State, state)
}

func GetESIOAuth2State(k app.Session) string {
	v, _ := k.Get(esiOAuth2State)
	if s, ok := v.(string); !ok {
		return ""
	} else {
		return s
	}
}

func ClearESIOAuth2State(k app.Session) {
	k.Delete(esiOAuth2State)
}

func SetCharacterSelected(k app.Session, cID int32) {
	k.Set(dharmaCharacterSelected, cID)
}

func GetCharacterSelected(k app.Session) int32 {
	v, _ := k.Get(dharmaCharacterSelected)
	if c, ok := v.(int32); !ok {
		return 0
	} else {
		return c
	}
}
