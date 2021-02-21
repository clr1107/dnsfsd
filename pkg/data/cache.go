package data

import (
	"sync"
	"time"
)

func now() int64 {
	return time.Now().Local().UnixNano() / 1000000
}

// ICache is an interface for Caches. Contains a standard list of methods all
// Caches should have. Note: not all caches are thread-safe.
//
// Put inserts a value of any type with a string key and the provided ttl (in
// milliseconds). Returns success.
//
// PutDefault calls #Put but with the default ttl for a given key; or the
// cache's default ttl is that could not be found. Returns success.
//
// Get returns the value associated with a key if it has not expired. If it is
// not present, or expired, nil is returned.
//
// Remove will remove the key from the cache if it is present.
//
// Contains will return if the key is contained in the cache.
//
// SetDefaultTTL will set the cache's default ttl for a specific key to the ttl
// given (in milliseconds).
//
// Size will return the number of entries in the cache. Not necessarily the
// number of non-expired entries, however.
type ICache interface {
	Put(key string, val interface{}, ttl int64) bool
	PutDefault(key string, val interface{}) bool
	Get(key string) interface{}
	Remove(key string) bool
	Contains(key string) bool
	SetDefaultTTL(key string, ttl int64)
	Clear()
	Size() int
}

// SimpleCache is a thread-safe implementation of ICache, simply uses a map with
// string keys, and a sync.RWMutex.
type SimpleCache struct {
	data        map[string]simpleCacheCell
	defaultTtls map[string]int64
	DefaultTTL  int64
	lock        *sync.RWMutex
}

// NewSimpleCache creates a new SimpleCache wtih the given default ttl.
func NewSimpleCache(defaultTTL int64) *SimpleCache {
	s := &SimpleCache{DefaultTTL: defaultTTL}

	s.data = make(map[string]simpleCacheCell)
	s.defaultTtls = make(map[string]int64)
	s.lock = &sync.RWMutex{}

	return s
}

type simpleCacheCell struct {
	val    interface{}
	expiry int64
}

func (c *simpleCacheCell) valid() bool {
	return c.expiry == -1 || now() < c.expiry
}

func (s *SimpleCache) read(key string) (simpleCacheCell, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	a, b := s.data[key]
	return a, b
}

func (s *SimpleCache) write(key string, cell *simpleCacheCell) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data[key] = *cell
}

func (s *SimpleCache) Put(key string, val interface{}, ttl int64) bool {
	var expiry int64

	if ttl == -1 {
		expiry = -1
	} else {
		expiry = now() + ttl
	}

	s.write(key, &simpleCacheCell{val, expiry})
	return true
}

func (s *SimpleCache) PutDefault(key string, val interface{}) bool {
	s.lock.RLock()
	ttl, ok := s.defaultTtls[key]
	s.lock.RUnlock()

	if !ok {
		ttl = s.DefaultTTL
	}

	return s.Put(key, val, ttl)
}

func (s *SimpleCache) Get(key string) interface{} {
	val, ok := s.read(key)

	if ok && !val.valid() {
		s.Remove(key)
		return nil
	}

	return val.val
}

func (s *SimpleCache) Remove(key string) bool {
	_, ok := s.read(key)

	if ok {
		s.lock.Lock()
		defer s.lock.Unlock()

		delete(s.data, key)
		return true
	}

	return false
}

func (s *SimpleCache) Contains(key string) bool {
	val, ok := s.read(key)

	if ok && !val.valid() {
		s.Remove(key)
		return false
	}

	return ok
}

func (s *SimpleCache) SetDefaultTtl(key string, ttl int64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.defaultTtls[key] = ttl
}

func (s *SimpleCache) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data = make(map[string]simpleCacheCell)
}

func (s *SimpleCache) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.data)
}
