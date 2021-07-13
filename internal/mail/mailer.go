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

package mail

import (
	"github.com/cjslep/dharma/internal/db"
	"github.com/go-fed/apcore/util"
)

type Mailer struct {
	DB *db.DB
}

func New(db *db.DB) *Mailer {
	return &Mailer{
		DB: db,
	}
}

func (m *Mailer) SendValidationEmail(c util.Context, userID string) error {
	// TODO: Make this more robust to partial failures
	err := m.DB.AddUserEmailValidationTask(c, userID)
	if err != nil {
		return err
	}
	// TODO: Send email
	return m.DB.MarkUserValidationEmailSent(c, userID)
}
