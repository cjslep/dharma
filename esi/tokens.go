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

package esi

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

type Tokens struct {
	Access        string
	Refresh       string
	AccessExpires time.Time
	CID           int
	CName         string
}

func NewTokens(jwt *JWTResponse, c *jwt.Claims) (*Tokens, error) {
	subs := strings.Split(c.Subject, ":")
	if len(subs) != 3 {
		return nil, fmt.Errorf("malformed subject in jwt: %s", c.Subject)
	}
	cid, err := strconv.Atoi(subs[2])
	if err != nil {
		return nil, fmt.Errorf("cannot convert character id to int: %w", err)
	}
	name, ok := c.Set["name"].(string)
	if !ok {
		return nil, fmt.Errorf("cannot fetch character name from set: %v", c.Set)
	}
	// TODO: Determine if other claims need to be processed
	return &Tokens{
		Access:        jwt.AccessToken,
		Refresh:       jwt.RefreshToken,
		AccessExpires: c.Expires.Time(),
		CID:           cid,
		CName:         name,
	}, nil
}
