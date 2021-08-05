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
	"context"

	"github.com/cjslep/dharma/internal/db"
	"github.com/cjslep/dharma/internal/mail"
	"github.com/go-fed/apcore/app"
	"golang.org/x/text/language"
)

type Users struct {
	F  app.Framework
	M  *mail.Mailer
	DB *db.DB
}

func InitAsCommandLineAdminUser(ctx context.Context, db *db.DB, userID string) error {
	// TODO: Generate token
	var token string
	err := db.AddUserEmailValidationTask(ctx, userID, token)
	if err != nil {
		return err
	}
	return db.MarkUserVerified(ctx, token)
}

func (u *Users) CreateUser(ctx context.Context, username, email, password string, lang language.Tag) error {
	// TODO: Make this more robust to partial failures
	userID, err := u.F.CreateUser(ctx, username, email, password)
	if err != nil {
		return err
	}
	// TODO: Generate token
	var token string
	err = u.DB.AddUserEmailValidationTask(ctx, userID, token)
	if err != nil {
		return err
	}
	if err = u.M.SendValidationEmail(ctx, userID, email, token, lang); err != nil {
		return err
	}
	return u.DB.MarkUserValidationEmailSent(ctx, userID)
}

func (u *Users) MarkUserVerified(c context.Context, token string) error {
	return u.DB.MarkUserVerified(c, token)
}

func (u *Users) IsUserVerified(c context.Context, userID string) (bool, error) {
	return u.DB.IsUserVerified(c, userID)
}
