package pokecache

import (
	"testing"
	"time"
)

func TestAddAndGet(t *testing.T) {
	cache := NewCache(5 * time.Second)
	cache.Add("https://pokeapi.co/api/v2/location-area?offset=0", []byte("test data"))

	val, ok := cache.Get("https://pokeapi.co/api/v2/location-area?offset=0")
	if !ok {
		t.Error("expected entry to be found in cache")
	}
	if string(val) != "test data" {
		t.Errorf("expected 'test data', got '%s'", string(val))
	}
}

func TestGetMiss(t *testing.T) {
	cache := NewCache(5 * time.Second)
	_, ok := cache.Get("nonexistent-key")
	if ok {
		t.Error("expected cache miss, got hit")
	}
}

func TestReap(t *testing.T) {
	interval := 50 * time.Millisecond
	cache := NewCache(interval)
	cache.Add("reap-me", []byte("goodbye"))

	time.Sleep(150 * time.Millisecond)

	_, ok := cache.Get("reap-me")
	if ok {
		t.Error("expected entry to be reaped, but it still exists")
	}
}