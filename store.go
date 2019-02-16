package main

import (
  "os"
  "io/ioutil"
  "log"
  "time"
)

// Store is just a very simple filesystem-hosted key-value store
type Store struct {
  Path string
}

// NewStore sets up a new Store, creating a directory if necessary
func NewStore(path string) *Store {
  store := &Store{Path: path}
  err := os.MkdirAll(path, 0700)
  if err != nil {
    log.Fatalf("unable to provision store directory %q: %v", path, err)
  }
  return store
}

// formatKey takes a bare key and translates it to a store filename
func (s *Store) formatKey(key string) string {
  return s.Path + "/" + key
}

// Get attempts to retrieve an item from the store. If this fails for
// any reason, or if the item is older than the specified TTL, nil and
// false are returned.
func (s *Store) Get(key string, ttl time.Duration) ([]byte,bool) {
  fkey := s.formatKey(key)
  stat,err := os.Stat(fkey)
  if err == nil {
    // matching key in cache
    now := time.Now()
    mtime := stat.ModTime()
    if now.Sub(mtime) > ttl {
      // but it's too old, invalidate it
      if err := os.Remove(fkey); err != nil {
        log.Printf("error purging cache for %q: %v", fkey, err)
      }
    } else {
      // unexpired cache key, try to read it
      if data,err := ioutil.ReadFile(fkey); err == nil {
        // cache key valid and unexpired
        return data,true
      }
    }
  }
  return nil, false
}

// Put simply stores an item in the Store
func (s *Store) Put(key string, data []byte) error {
  return ioutil.WriteFile(s.formatKey(key), data, 0700)
}
