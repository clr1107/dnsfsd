package cache

import (
	"strconv"
	"testing"
)

func TestPutGet(t *testing.T) {
	cache := NewSimpleCache(-1)

	cache.PutDefault("one", int64(1))
	one := cache.Get("one")

	if one == nil {
		t.Fatalf("value was nil (not in the cache)")
	}

	if val, ok := one.(int64); !ok {
		t.Fatalf("value given could not be cast to int64")
	} else {
		if val != 1 {
			t.Fatalf("not the correct value (not int64(1))")
		}
	}
}

func TestTimeout(t *testing.T) {
	cache := NewSimpleCache(-2)

	cache.PutDefault("one", 1)
	one := cache.Get("one")

	if one != nil {
		t.Fatalf("value is still being returned despite it expiring")
	}
}

func TestContains(t *testing.T) {
	cache := NewSimpleCache(-1)
	cache.PutDefault("one", 1)

	if !cache.Contains("one") {
		t.Fatalf("cache does not contain key just added")
	}

	if !cache.Remove("one") {
		t.Fatalf("could not remove key added")
	}

	if cache.Contains("one") {
		t.Fatalf("cache still contains key just removed")
	}
}

func TestSize(t *testing.T) {
	cache := NewSimpleCache(-1)

	for i := 0; i < 100; i++ {
		cache.PutDefault(strconv.Itoa(i), i)
	}

	size := cache.Size()
	if size != 100 {
		t.Fatalf("size of cache should be 100, but it is %v", size)
	}

	cache.Clear()
	size = cache.Size()

	if size != 0 {
		t.Fatalf("cache was cleared and size is %v, not 0", size)
	}
}

func TestSerialisation(t *testing.T) {
	cache := NewDNSCache(-1)

	cache.PutDefault("one", 1)
	cache.PutDefault("two", 2)
	cache.PutDefault("three", 3)

	buffer, err := cache.Serialise()
	if err != nil {
		t.Fatalf("%v", err)
	}

	deserialised, err := DeserialiseDNSCache(buffer)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if deserialised.Size() != cache.Size() {
		t.Fatalf("deserialised's size is not equal to the original")
	}

	for k, v := range deserialised.Data {
		if v.Val != cache.Data[k].Val {
			t.Fatalf("one or more elements are not correct. %v != %v", v.Val, cache.Data[k].Val)
		}
	}
}