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
	"github.com/cjslep/dharma/internal/mail"
	"github.com/go-fed/apcore/app"
	"github.com/go-fed/apcore/util"
)

type Users struct {
	F  app.Framework
	M  *mail.Mailer
	DB *db.DB
}

func (u *Users) CreateUser(ctx util.Context, username, email, password string) error {
	// TODO: Make this more robust to partial failures
	userID, err := u.F.CreateUser(ctx, username, email, password)
	if err != nil {
		return err
	}
	err = u.M.SendValidationEmail(ctx, userID)
	return err
}

func (u *Users) MarkUserValidated(c util.Context, userID string) error {
	return u.DB.MarkUserValidated(c, userID)
}

func (u *Users) IsUserValidated(c util.Context, userID string) (bool, error) {
	return u.DB.IsUserValidated(c, userID)
}

func (u *Users) HasEveAccount() {
	// TODO
}

func (u *Users) HasCharacterSelected() {
	// TODO
}
