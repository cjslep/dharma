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

package account

import (
	"github.com/cjslep/dharma/internal/api"
	"github.com/go-fed/apcore/app"
)

type Account struct {
	C *api.Context
}

func (a *Account) Route(r app.Router) {
	r.NewRoute().Methods("GET").WebOnlyHandlerFunc("/account/profile", api.MustHaveSessionAndLanguageCode(a.C, a.getProfile))
	r.NewRoute().Methods("POST").WebOnlyHandlerFunc("/account/profile", api.MustHaveSessionAndLanguageCode(a.C, a.postProfile))
	r.NewRoute().Methods("GET").WebOnlyHandlerFunc("/account/settings", api.MustHaveSessionAndLanguageCode(a.C, a.getSettings))
	r.NewRoute().Methods("POST").WebOnlyHandlerFunc("/account/settings", api.MustHaveSessionAndLanguageCode(a.C, a.postSettings))
	r.NewRoute().Methods("GET").WebOnlyHandlerFunc("/account/characters", api.MustHaveSessionAndLanguageCode(a.C, a.getCharacters))
	r.NewRoute().Methods("POST").WebOnlyHandlerFunc("/account/characters", api.MustHaveSessionAndLanguageCode(a.C, a.postCharacters))
}
