package persistence

import (
	"sync"
	"time"
)

func now() int64 {
	return time.Now().Local().UnixNano() / 1000000
}

type ICache interface {
	Put(key string, val interface{}, ttl int64) bool
	PutDefault(key string, val interface{}) bool
	Get(key string) interface{}
	Remove(key string) bool
	Contains(key string) bool
	SetDefaultTtl(key string, ttl int64)
	Clear()
	Size() int
}

type SimpleCache struct {
	data        map[string]simpleCacheCell
	defaultTtls map[string]int64
	DefaultTtl  int64
	lock        *sync.RWMutex
}

func NewSimpleCache(defaultTtl int64) *SimpleCache {
	s := &SimpleCache{DefaultTtl: defaultTtl}

	s.data = make(map[string]simpleCacheCell)
	s.defaultTtls = make(map[string]int64)
	s.lock = &sync.RWMutex{}

	return s
}

type simpleCacheCell struct {
	val    interface{}
	expiry int64
}

func (c *simpleCacheCell) Valid() bool {
	return c.expiry == -1 || now() < c.expiry
}

func (s *SimpleCache) Put(key string, val interface{}, ttl int64) bool {
	var expiry int64

	if ttl == -1 {
		expiry = -1
	} else {
		expiry = now() + ttl
	}

	s.lock.Lock()
	s.data[key] = simpleCacheCell{val, expiry}
	s.lock.Unlock()

	return true
}

func (s *SimpleCache) PutDefault(key string, val interface{}) bool {
	if ttl, ok := s.defaultTtls[key]; ok {
		return s.Put(key, val, ttl)
	}

	return s.Put(key, val, s.DefaultTtl)
}

func (s *SimpleCache) Get(key string) interface{} {
	if !s.Contains(key) {
		return nil
	}

	s.lock.RLock()
	val := s.data[key]
	s.lock.RUnlock()

	if !val.Valid() {
		s.Remove(key)
		return nil
	}

	return val.val
}

func (s *SimpleCache) Remove(key string) bool {
	if s.Contains(key) {
		s.lock.Lock()

		delete(s.data, key)

		s.lock.Unlock()
		return true
	}

	return false
}

func (s *SimpleCache) Contains(key string) bool {
	s.lock.RLock()
	_, ok := s.data[key]
	s.lock.RUnlock()

	return ok
}

func (s *SimpleCache) SetDefaultTtl(key string, ttl int64) {
	s.defaultTtls[key] = ttl
}

func (s *SimpleCache) Clear() {
	s.lock.Lock()
	s.data = make(map[string]simpleCacheCell)
	s.lock.Unlock()

	s.lock = &sync.RWMutex{}
}

func (s *SimpleCache) Size() int {
	s.lock.RLock()
	l := len(s.data)
	s.lock.RUnlock()

	return l
}
