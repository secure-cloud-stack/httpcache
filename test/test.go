package test

import (
	"bytes"
	"testing"

	"github.com/secure-cloud-stack/httpcache"
)

// Cache excercises a httpcache.Cache implementation.
func Cache(t *testing.T, cache httpcache.Cache) {
	key := "testKey"
	_, ok, _ := cache.Get(key)
	if ok {
		t.Fatal("retrieved key before adding it")
	}

	val := []byte("some bytes")
	_ = cache.Set(key, val)

	retVal, ok, _ := cache.Get(key)
	if !ok {
		t.Fatal("could not retrieve an element we just added")
	}
	if !bytes.Equal(retVal, val) {
		t.Fatal("retrieved a different value than what we put in")
	}

	_ = cache.Delete(key)

	_, ok, _ = cache.Get(key)
	if ok {
		t.Fatal("deleted key still present")
	}
}
