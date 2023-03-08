//go:build !appengine
// +build !appengine

// Package memcache provides an implementation of httpcache.Cache that uses
// gomemcache to store cached responses.
//
// When built for Google App Engine, this package will provide an
// implementation that uses App Engine's memcache service.  See the
// appengine.go file in this package for details.
package memcache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/secure-cloud-stack/httpcache"
)

// Cache is an implementation of httpcache.Cache that caches responses in a
// memcache server.
type Cache struct {
	*memcache.Client
}

// cacheKey modifies an httpcache key for use in memcache.  Specifically, it
// prefixes keys to avoid collision with other data stored in memcache.
func cacheKey(key string) string {
	return "httpcache:" + key
}

// Get returns the response corresponding to key if present.
func (c *Cache) Get(key string) (resp []byte, ok bool, err error) {
	item, err := c.Client.Get(cacheKey(key))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, false, nil
		}
		return nil, false, err
	}
	return item.Value, true, nil
}

// Set saves a response to the cache as key.
func (c *Cache) Set(key string, resp []byte) error {
	item := &memcache.Item{
		Key:   cacheKey(key),
		Value: resp,
	}
	return c.Client.Set(item)
}

// Delete removes the response with key from the cache.
func (c *Cache) Delete(key string) error {
	return c.Client.Delete(cacheKey(key))
}

// New returns a new Cache using the provided memcache server(s) with equal
// weight. If a server is listed multiple times, it gets a proportional amount
// of weight.
func New(server ...string) httpcache.Cache {
	return NewWithClient(memcache.New(server...))
}

// NewWithClient returns a new Cache with the given memcache client.
func NewWithClient(client *memcache.Client) httpcache.Cache {
	return &Cache{client}
}
