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

package site

import (
	"github.com/cjslep/dharma/internal/api"
	"github.com/go-fed/apcore/app"
)

type Site struct {
	C *api.Context
}

func (s *Site) Route(r app.Router) {
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/",
		api.CorpMustBeManaged(s.C,
			api.MustHaveLanguageCode(s.getHome)))
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/about",
		api.CorpMustBeManaged(s.C,
			api.MustHaveLanguageCode(s.getAbout)))
	r.NewRoute().Methods("GET").WebOnlyHandler(
		"/site/setup/corp",
		api.CorpMustNotBeManaged(s.C,
			api.MustBeAdmin(s.C,
				api.MustHaveLanguageCode(s.getChooseCorpToManage))))
	r.NewRoute().Methods("POST").WebOnlyHandler(
		"/site/setup/corp",
		api.CorpMustNotBeManaged(s.C,
			api.MustBeAdmin(s.C,
				api.MustHaveSessionAndLanguageCode(s.C, s.postChooseCorpToManage))))
	r.NewRoute().Methods("POST").WebOnlyHandler(
		"/site/setup/corp/search",
		api.CorpMustNotBeManaged(s.C,
			api.MustBeAdmin(s.C,
				api.MustHaveLanguageCode(s.postCorpSetupSearch))))
}
