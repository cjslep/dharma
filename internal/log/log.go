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

package log

import (
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Logger(debug bool, dir, filename string, backups, size, age int) *zerolog.Logger {
	if size < 0 {
		size = 100
	}
	if backups < 0 {
		backups = 0
	}
	if age < 0 {
		age = 0
	}
	writers := make([]io.Writer, 0, 2)
	writers = append(writers, &lumberjack.Logger{
		Filename:   path.Join(dir, filename),
		MaxBackups: backups,
		MaxSize:    size,
		MaxAge:     age,
	})
	if debug {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	w := io.MultiWriter(writers...)
	l := zerolog.New(w).With().Timestamp().Logger()
	return &l
}
