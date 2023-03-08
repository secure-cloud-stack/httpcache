// Package diskcache provides an implementation of httpcache.Cache that uses the diskv package
// to supplement an in-memory map with persistent storage
package diskcache

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/peterbourgon/diskv"
	"github.com/secure-cloud-stack/httpcache"
)

// Cache is an implementation of httpcache.Cache that supplements the in-memory map with persistent storage
type Cache struct {
	d *diskv.Diskv
}

// Get returns the response corresponding to key if present
func (c *Cache) Get(key string) (resp []byte, ok bool, err error) {
	key = keyToFilename(key)
	resp, err = c.d.Read(key)
	if err != nil {
		return []byte{}, false, err
	}
	return resp, true, nil
}

// Set saves a response to the cache as key
func (c *Cache) Set(key string, resp []byte) error {
	key = keyToFilename(key)
	return c.d.WriteStream(key, bytes.NewReader(resp), true)
}

// Delete removes the response with key from the cache
func (c *Cache) Delete(key string) error {
	key = keyToFilename(key)
	return c.d.Erase(key)
}

func keyToFilename(key string) string {
	s := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", s)
}

// New returns a new Cache that will store files in basePath
func New(basePath string) *Cache {
	return &Cache{
		d: diskv.New(diskv.Options{
			BasePath:     basePath,
			CacheSizeMax: 100 * 1024 * 1024, // 100MB
		}),
	}
}

// NewWithDiskv returns a new Cache using the provided Diskv as underlying
// storage.
func NewWithDiskv(d *diskv.Diskv) httpcache.Cache {
	return &Cache{d}
}
