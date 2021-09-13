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

package util

import (
	"bytes"
	"image"
	"image/jpeg"

	_ "image/gif"
	_ "image/png"
)

// toImage turns bytes into a JPEG, GIF, or PNG image.
func ToImage(b []byte) (image.Image, error) {
	buf := bytes.NewBuffer(b)
	i /*image type=*/, _, err := image.Decode(buf)
	return i, err
}

func EncodeJPEG(i image.Image) ([]byte, error) {
	var b bytes.Buffer
	o := &jpeg.Options{Quality: 100}
	err := jpeg.Encode(&b, i, o)
	return b.Bytes(), err
}
