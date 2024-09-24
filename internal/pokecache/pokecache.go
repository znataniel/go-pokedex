package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	stored   map[string]cacheEntry
	halfLife time.Duration
	mu       *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(duration int) Cache {
	newCache := Cache{
		stored:   make(map[string]cacheEntry),
		halfLife: time.Duration(duration * int(time.Second)),
		mu:       &sync.Mutex{},
	}

	go newCache.reapLoop()

	return newCache
}

func (cache *Cache) Add(key string, val []byte) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if _, exists := cache.stored[key]; exists {
		return fmt.Errorf("cache add error: entry already exists")
	}

	cache.stored[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	return nil
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	entry, exists := cache.stored[key]
	if !exists {
		return nil, false
	}

	return entry.val, true
}

func (cache *Cache) reapLoop() {
	ticker := time.NewTicker(cache.halfLife)

	for {
		select {
		case <-ticker.C:
			cache.mu.Lock()
			defer cache.mu.Unlock()

			for key, value := range cache.stored {
				if time.Now().After(value.createdAt.Add(cache.halfLife)) {
					delete(cache.stored, key)
				}
			}
		}
	}
}
