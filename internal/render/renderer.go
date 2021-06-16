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

package render

import (
	"github.com/cjslep/dharma/locales"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

type Renderer struct {
	b *i18n.Bundle
}

func New() (*Renderer, error) {
	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	if err := locales.AddMessageFiles(b); err != nil {
		return nil, err
	}

	return &Renderer{
		b: b,
	}, nil
}
