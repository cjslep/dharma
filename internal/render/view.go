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
	"io"
)

type View struct {
	w        io.Writer
	status   int
	htmlData *htmlData
	jsonData *jsonData
}

type htmlData struct {
	Name string
	Data map[string]interface{}
}

type jsonData struct {
	Payload interface{}
}

func NewHTMLView(w io.Writer, status int, name string, data map[string]interface{}) *View {
	return &View{
		w:      w,
		status: status,
		htmlData: &htmlData{
			Name: name,
			Data: data,
		},
	}
}

func NewJSONView(w io.Writer, status int, payload interface{}) *View {
	return &View{
		w:      w,
		status: status,
		jsonData: &jsonData{
			Payload: payload,
		},
	}
}

func (v *View) isHTML() bool {
	return v.htmlData != nil
}

func (v *View) isJSON() bool {
	return v.jsonData != nil
}
