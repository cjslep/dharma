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

// +build dev

package assets

import (
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	prefix0 = "assets"
	prefix1 = "src"
)

var assets http.FileSystem = http.Dir(filepath.Join(prefix0, prefix1))

func Asset(name string) ([]byte, error) {
	bs, err := ioutil.ReadFile(filepath.Join(prefix0, prefix1, name))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return bs, nil
}

func AssetNames() []string {
	var out []string
	dir := filepath.Join(prefix0, prefix1)
	err := filepath.WalkDir(dir,
		func(path string, d fs.DirEntry, err error) error {
			locs := strings.Split(path, dir)
			if len(locs) != 2 {
				return errors.New("cannot walk asset tree: split could not get sub-file name")
			}
			if d.IsDir() {
				return nil
			}
			name := locs[1][1:] // Omit leading slash
			if !strings.HasSuffix(name, ".tmpl") {
				return nil
			}
			out = append(out, name)
			return nil
		})
	if err != nil {
		panic(err) // Only permitted with +build dev
	}
	return out
}
