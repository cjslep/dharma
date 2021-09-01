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
	"github.com/cjslep/dharma/internal/api/paths"
	"github.com/go-fed/apcore/app"
)

type Account struct {
	C *api.Context
}

func (a *Account) Route(r app.Router) {
	r.NewRoute().Methods("GET").WebOnlyHandler(
		paths.VerifyPath,
		api.MustHaveLanguageCode(a.getVerify)) // TODO: Make POST handler as well
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/account/register",
		api.MustHaveLanguageCode(a.getRegister))
	r.NewRoute().Methods("POST").WebOnlyHandler(
		"/account/register",
		api.MustHaveLanguageCode(a.postRegister))
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/account/profile",
		api.CorpMustBeManaged(a.C,
			api.MustHaveLanguageCode(a.getProfile)))
	r.NewRoute().Methods("POST").WebOnlyHandler(
		"/account/profile",
		api.CorpMustBeManaged(a.C,
			api.MustHaveLanguageCode(a.postProfile)))
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/account/settings",
		api.CorpMustBeManaged(a.C,
			api.MustHaveLanguageCode(a.getSettings)))
	r.NewRoute().Methods("POST").WebOnlyHandler(
		"/account/settings",
		api.CorpMustBeManaged(a.C,
			api.MustHaveLanguageCode(a.postSettings)))
	r.NewRoute().Methods("GET").WebOnlyHandler(
		paths.AccountCharactersPath,
		api.MustHaveSessionAndLanguageCode(a.C, a.getCharacters))
	r.NewRoute().Methods("POST").WebOnlyHandler(
		"/account/characters",
		api.MustHaveSessionAndLanguageCode(a.C, a.postCharacters))
}
