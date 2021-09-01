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
	"net/http"

	"github.com/cjslep/dharma/internal/api"
	"github.com/cjslep/dharma/internal/render"
	"github.com/cjslep/dharma/internal/sessions"
	"github.com/go-fed/apcore/app"
	"github.com/go-fed/apcore/util"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

func (a *Account) getCharacters(w http.ResponseWriter, r *http.Request, k app.Session, langs []language.Tag) {
	userID, err := k.UserID()
	if err != nil {
		a.C.MustRenderError(w, r, errors.Wrap(err, "error getting user for session"), langs...)
		return
	}

	chars, err := a.C.ESI.GetCharactersForUser(util.Context{r.Context()}, userID)
	if err != nil {
		a.C.MustRenderError(w, r, errors.Wrap(err, "error getting characters for session"), langs...)
		return
	}

	rc := api.From(r.Context())
	v := render.NewHTMLView(
		w,
		http.StatusOK,
		"account/characters",
		rc,
		map[string]interface{}{
			"characters": chars,
			"selectedID": sessions.GetCharacterSelected(k),
		},
		langs...)
	a.C.MustRender(v)
}
