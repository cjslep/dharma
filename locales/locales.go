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

package locales

import (
	"io/ioutil"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
)

var messageFiles = []string{
	"active.de-DE.toml",
	"active.en.toml",
	"active.ru-RU.toml",
	"active.zh-Hans.toml",
	"active.zh-Hant.toml",
}

func AddMessageFiles(b *i18n.Bundle) error {
	for _, f := range messageFiles {
		file, err := assets.Open(f)
		if err != nil {
			return errors.Wrap(err, "could not open localized message file: "+f)
		}
		buf, err := ioutil.ReadAll(file)
		if err != nil {
			return errors.Wrap(err, "error reading localized message file: "+f)
		}
		_, err = b.ParseMessageFileBytes(buf, f)
		if err != nil {
			return errors.Wrap(err, "could not parse localized message file: "+f)
		}
	}
	return nil
}
