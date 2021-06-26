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

// +build tools

package tools

// Below, we list the dependencies needed at `go generate` time. Running the
// generate command requires modules not normally needed in a regular build of
// dharma. Therefore, `go generate` demands `go get ...` for the below
// dependencies, but `go mod tidy` "helpfully" removes those modules, creating
// thrashing in the go.mod file.
//
// This is a known problem to the go folks:
//
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://github.com/golang/go/issues/25922
// https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md
//
// The convention is to use the `tools` build constraint instead of the `ignore`
// one, in an act to ensure it is really and truly ignored.

import (
	"github.com/shurcooL/vfsgen"
)
