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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/pascaldekloe/jwt"
	"github.com/pkg/errors"
)

func FetchEveOnlineKeys() (*OAuthKeysMetadata, error) {
	resp, err := http.Get("https://login.eveonline.com/oauth/jwks")
	if err != nil {
		return nil, fmt.Errorf("error fetching ESI JWT public keys: %w", err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading public key body: %w", err)
		} else {
			k := &OAuthKeysMetadata{}
			err := json.Unmarshal(body, k)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling public key body: %w", err)
			} else {
				return k, nil
			}
		}
	}
}

type OAuthKeysMetadata struct {
	Keys []*KeysMetadata `json:"keys"`
}

var _ driver.Valuer = OAuthKeysMetadata{}
var _ sql.Scanner = &OAuthKeysMetadata{}

func (o *OAuthKeysMetadata) JWTKey() *KeysMetadata {
	for _, v := range o.Keys {
		if v.Id == "JWT-Signature-Key" {
			return v
		}
	}
	return nil
}

func (o OAuthKeysMetadata) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func (o *OAuthKeysMetadata) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("failed to assert scan src to []byte type")
	}
	return json.Unmarshal(b, o)
}

type KeysMetadata struct {
	Alg string `json:"alg"`
	Id  string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	// RS
	N string `json:"n"`
	E string `json:"e"`
	// ES
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

func (k *KeysMetadata) ValidateToken(b []byte) (*jwt.Claims, error) {
	if k.isRSAKey() {
		r := k.rsaKey()
		return jwt.RSACheck(b, r)
	} else if k.isECDSAKey() {
		e := k.ecdsaKey()
		return jwt.ECDSACheck(b, e)
	} else {
		return nil, fmt.Errorf("unhandled jwt key type: %s", k.Alg)
	}
}

func (k *KeysMetadata) isRSAKey() bool {
	_, ok := jwt.RSAAlgs[k.Alg]
	return ok
}

func (k *KeysMetadata) isECDSAKey() bool {
	_, ok := jwt.ECDSAAlgs[k.Alg]
	return ok
}

func (k *KeysMetadata) rsaKey() *rsa.PublicKey {
	bn, _ := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(k.N)
	be, _ := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(k.E)
	n := big.NewInt(0)
	n = n.SetBytes(bn)
	e := big.NewInt(0)
	e = e.SetBytes(be)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}
}

func (k *KeysMetadata) ecdsaKey() *ecdsa.PublicKey {
	bx, _ := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(k.X)
	by, _ := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(k.Y)
	x := big.NewInt(0)
	x = x.SetBytes(bx)
	y := big.NewInt(0)
	y = y.SetBytes(by)

	var crv elliptic.Curve
	switch k.Crv {
	case "P-224":
		crv = elliptic.P224()
	case "P-256":
		crv = elliptic.P256()
	case "P-384":
		crv = elliptic.P384()
	case "P-521":
		crv = elliptic.P521()
	}
	return &ecdsa.PublicKey{
		Curve: crv,
		X:     x,
		Y:     y,
	}
}
