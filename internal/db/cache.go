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

package db

import (
	"sync"
)

type stringCache struct {
	v  string
	mu sync.RWMutex
}

func (c *stringCache) Get() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v
}

func (c *stringCache) Set(v string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v = v
}

type cache struct {
	corpName *stringCache
}

func newCache() *cache {
	return &cache{
		corpName: &stringCache{},
	}
}

func (d *DB) GetCorpName() string {
	return d.cache.corpName.Get()
}

// TODO: Set corp name, at init and on the fly
func (d *DB) setCorpName(v string) {
	d.cache.corpName.Set(v)
}
