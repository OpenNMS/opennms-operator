/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package yaml

import (
	"sync"
)

var (
	yamlcache *cache
)

// cache - an in memory cache
type cache struct {
	yamlmap map[string]string
	mutex   sync.RWMutex
}

// Set - set a value for a given key in the cache
func (c *cache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.yamlmap[key] = value
}

// Get - get the value of a given key from the cache
func (c *cache) Get(key string) (value string, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item, ok := c.yamlmap[key]
	return item, ok
}

// Cache - init a new cache or return the existing one
func Cache() *cache {
	if yamlcache == nil {
		yamlcache = &cache{
			yamlmap: make(map[string]string),
		}
	}
	return yamlcache
}