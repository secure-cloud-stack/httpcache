// Package leveldbcache provides an implementation of httpcache.Cache that
// uses github.com/syndtr/goleveldb/leveldb
package leveldbcache

import (
	"github.com/secure-cloud-stack/httpcache"
	"github.com/syndtr/goleveldb/leveldb"
)

// Cache is an implementation of httpcache.Cache with leveldb storage
type Cache struct {
	db *leveldb.DB
}

// Get returns the response corresponding to key if present
func (c *Cache) Get(key string) (resp []byte, ok bool, err error) {
	resp, err = c.db.Get([]byte(key), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return []byte{}, false, nil
		}
		return []byte{}, false, err
	}
	return resp, true, nil
}

// Set saves a response to the cache as key
func (c *Cache) Set(key string, resp []byte) error {
	return c.db.Put([]byte(key), resp, nil)
}

// Delete removes the response with key from the cache
func (c *Cache) Delete(key string) error {
	return c.db.Delete([]byte(key), nil)
}

// New returns a new Cache that will store leveldb in path
func New(path string) (*Cache, error) {
	cache := &Cache{}

	var err error
	cache.db, err = leveldb.OpenFile(path, nil)

	if err != nil {
		return nil, err
	}
	return cache, nil
}

// NewWithDB returns a new Cache using the provided leveldb as underlying
// storage.
func NewWithDB(db *leveldb.DB) httpcache.Cache {
	return &Cache{db}
}
