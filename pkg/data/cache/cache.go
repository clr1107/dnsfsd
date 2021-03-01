package cache

import (
	"encoding/gob"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"os"
	"sync"
	"time"
)

// ICache is an interface for Caches. Contains a standard list of methods all
// Caches should have. Note: not all caches are thread-safe.
//
// Put inserts a value of any type with a string key and the provided ttl.
// Returns success.
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
//
// This is just a wrapper around whatever cache library is being used -- to
// allow quick changes.
type ICache interface {
	Put(key string, val interface{}, ttl time.Duration) bool
	PutDefault(key string, val interface{}) bool
	Get(key string) interface{}
	Remove(key string) bool
	Contains(key string) bool
	SetDefaultTTL(key string, ttl time.Duration)
	Clear()
	Clean()
	Size() int
}

// SimpleCache is a thread-safe implementation of ICache, simply uses a map with
// string keys, and a sync.RWMutex.
type SimpleCache struct {
	Impl        *cache.Cache
	DefaultTTLs map[string]time.Duration
	DefaultTTL  time.Duration
	ttlLock     *sync.RWMutex
}

// NewSimpleCache creates a new SimpleCache with the given default ttl.
func NewSimpleCache(defaultTTL time.Duration) *SimpleCache {
	s := &SimpleCache{DefaultTTL: defaultTTL}

	s.Impl = cache.New(defaultTTL, 5*time.Minute)
	s.DefaultTTLs = make(map[string]time.Duration)
	s.ttlLock = &sync.RWMutex{}

	return s
}

func (s *SimpleCache) Put(key string, val interface{}, ttl time.Duration) bool {
	s.Impl.Set(key, val, ttl)
	return true
}

func (s *SimpleCache) PutDefault(key string, val interface{}) bool {
	s.ttlLock.RLock()
	ttl, ok := s.DefaultTTLs[key]
	s.ttlLock.RUnlock()

	if !ok {
		ttl = s.DefaultTTL
	}

	return s.Put(key, val, ttl)
}

func (s *SimpleCache) Get(key string) interface{} {
	val, ok := s.Impl.Get(key)

	if !ok {
		return nil
	} else {
		return val
	}
}

func (s *SimpleCache) Remove(key string) bool {
	ok := s.Contains(key)

	if ok {
		s.Impl.Delete(key)
	}

	return ok
}

func (s *SimpleCache) Contains(key string) bool {
	_, ok := s.Impl.Get(key)
	return ok
}

func (s *SimpleCache) SetDefaultTtl(key string, ttl time.Duration) {
	s.ttlLock.Lock()
	s.DefaultTTLs[key] = ttl
	s.ttlLock.Unlock()
}

func (s *SimpleCache) Clear() {
	s.Impl.Flush()
}

func (s *SimpleCache) Clean() {
	s.Impl.DeleteExpired()
}

func (s *SimpleCache) Size() int {
	return len(s.Impl.Items())
}

type DNSCache struct {
	*SimpleCache
}

func NewDNSCache(defaultTTL time.Duration) *DNSCache {
	return &DNSCache{NewSimpleCache(defaultTTL)}
}

func registerGobTypes() {
	gob.Register(dns.Question{})
	gob.Register(make([]dns.RR, 0))
	gob.Register(new(dns.A))
}

func DNSCacheFromFile(defaultTTL time.Duration, path string) (*DNSCache, error) {
	c := cache.New(defaultTTL, 5*time.Minute)
	fp, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	registerGobTypes()
	if err := c.Load(fp); err != nil {
		return nil, err
	}

	dc := NewDNSCache(defaultTTL)
	dc.Impl = c

	return dc, nil
}

func (d *DNSCache) SerialiseToFile(path string) error {
	registerGobTypes()
	return d.Impl.SaveFile(path)
}
