package cache

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"sync"
	"time"

	"github.com/miekg/dns"
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
	Put(key interface{}, val interface{}, ttl time.Duration) bool
	PutDefault(key interface{}, val interface{}) bool
	Get(key interface{}) interface{}
	Remove(key interface{}) bool
	Contains(key interface{}) bool
	SetDefaultTTL(key interface{}, ttl time.Duration)
	Clear()
	Clean()
	Size() int
}

// SimpleCache is a thread-safe implementation of ICache, simply uses a map with
// string keys, and a sync.RWMutex.
type SimpleCache struct {
	Data        map[interface{}]simpleCacheCell
	defaultTTLs map[interface{}]time.Duration
	DefaultTTL  time.Duration
	lock        *sync.RWMutex
}

// NewSimpleCache creates a new SimpleCache with the given default ttl.
func NewSimpleCache(defaultTTL time.Duration) *SimpleCache {
	s := &SimpleCache{DefaultTTL: defaultTTL}

	s.Data = make(map[interface{}]simpleCacheCell)
	s.defaultTTLs = make(map[interface{}]time.Duration)
	s.lock = &sync.RWMutex{}

	return s
}

type simpleCacheCell struct {
	Val    interface{}
	Expiry *time.Time
}

func (c *simpleCacheCell) valid() bool {
	return c.Expiry == nil || c.Expiry.After(time.Now())
}

func (s *SimpleCache) read(key interface{}) (simpleCacheCell, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	a, b := s.Data[key]
	return a, b
}

func (s *SimpleCache) write(key interface{}, cell *simpleCacheCell) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Data[key] = *cell
}

func (s *SimpleCache) Put(key interface{}, val interface{}, ttl time.Duration) bool {
	var expiry *time.Time = nil

	if ttl.Seconds() != -1 {
		x := time.Now().Add(ttl)
		expiry = &x
	}

	s.write(key, &simpleCacheCell{val, expiry})
	return true
}

func (s *SimpleCache) PutDefault(key interface{}, val interface{}) bool {
	s.lock.RLock()
	ttl, ok := s.defaultTTLs[key]
	s.lock.RUnlock()

	if !ok {
		ttl = s.DefaultTTL
	}

	return s.Put(key, val, ttl)
}

func (s *SimpleCache) Get(key interface{}) interface{} {
	val, ok := s.read(key)

	if ok && !val.valid() {
		s.Remove(key)
		return nil
	}

	return val.Val
}

func (s *SimpleCache) Remove(key interface{}) bool {
	_, ok := s.read(key)

	if ok {
		s.lock.Lock()
		defer s.lock.Unlock()

		delete(s.Data, key)
		return true
	}

	return false
}

func (s *SimpleCache) Contains(key interface{}) bool {
	s.lock.Lock()
	s.lock.RLock()

	defer s.lock.Unlock()
	defer s.lock.RUnlock()

	val, ok := s.read(key)

	if ok && !val.valid() {
		delete(s.Data, key)
		return false
	}

	return ok
}

func (s *SimpleCache) SetDefaultTtl(key interface{}, ttl time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.defaultTTLs[key] = ttl
}

func (s *SimpleCache) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Data = make(map[interface{}]simpleCacheCell)
}

func (s *SimpleCache) Clean() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for k, v := range s.Data {
		if !v.valid() {
			delete(s.Data, k)
		}
	}
}

func (s *SimpleCache) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.Data)
}

func registerGobTypes() {
	gob.Register(dns.Question{})
	gob.Register(make([]dns.RR, 0))
	gob.Register(new(dns.A))
}

type DNSCache struct {
	*SimpleCache
}

func NewDNSCache(defaultTTL time.Duration) *DNSCache {
	return &DNSCache{NewSimpleCache(defaultTTL)}
}

func DeserialiseDNSCache(buffer *bytes.Buffer) (*DNSCache, error) {
	d := new(DNSCache)

	registerGobTypes()
	decoder := gob.NewDecoder(buffer)

	if err := decoder.Decode(d); err != nil {
		return nil, err
	}

	d.lock = new(sync.RWMutex)
	return d, nil
}

func DNSCacheFromFile(path string) (*DNSCache, error) {
	read, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(read)
	return DeserialiseDNSCache(buffer)
}

func (d *DNSCache) Serialise() (*bytes.Buffer, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	buffer := new(bytes.Buffer)

	registerGobTypes()
	encoder := gob.NewEncoder(buffer)

	if err := encoder.Encode(d); err != nil {
		return nil, err
	}

	return buffer, nil
}

func (d *DNSCache) SerialiseToFile(path string) error {
	buffer, err := d.Serialise()

	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, buffer.Bytes(), 0666)
}
