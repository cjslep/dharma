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

package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Messages struct {
	l *i18n.Localizer
}

func New(b *i18n.Bundle, langs ...string) *Messages {
	return &Messages{
		l: i18n.NewLocalizer(b, langs...),
	}
}

func (m *Messages) PageNotFound() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "pageNotFound",
			Description: "Notifies the user that the webpage was not found",
			Other:       "The requested page was not found.",
		},
	})
}

func (m *Messages) BadRequest() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "badRequest",
			Description: "Notifies the user that the browser issued a bad request",
			Other:       "Your browser issued a bad request.",
		},
	})
}

func (m *Messages) MethodNotAllowed() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "methodNotAllowed",
			Description: "Notifies the user that the browser issued an unsupported HTTP method",
			Other:       "Your browser issued a request with a method that is disallowed",
		},
	})
}

func (m *Messages) UhOh() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "uhOh",
			Description: "Exclaimation that a mistake or accident happened",
			Other:       "Uh oh!",
		},
	})
}

func (m *Messages) InternalServerError() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "internalServerError",
			Description: "Notifies the user that the server had an internal error",
			Other:       "An unrecoverable error occurred within the server.",
		},
	})
}

func (m *Messages) UedamaGateCampRequest() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "uedamaGateCampRequest",
			Description: "Jovial comment that perhaps a user's HTTP request was gate-camped and blown up in an EVE Online system named Uedama",
			Other:       "Perhaps your request was blown up in an Uedama gate camp?",
		},
	})
}

func (m *Messages) ConsiderFileABug() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "fileABug",
			Description: "Prompts the user to click on this text in order to file an issue against the software",
			Other:       "Please consider filing a bug report",
		},
	})
}
