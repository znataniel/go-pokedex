package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const EXPTIME = time.Millisecond * 10
	const WAITTIME = time.Millisecond * 30

	cache := NewCache(EXPTIME)
	if err := cache.Add("https://test.case", []byte("testest")); err != nil {
		t.Errorf("adding entry to cache failed: %v", err)
	}

	if _, exists := cache.Get("https://test.case"); !exists {
		t.Errorf("entry not found immediately after adding it")
	}

	time.Sleep(WAITTIME)

	if _, exists := cache.Get("https://test.case"); exists {
		t.Errorf("entry still exists after reap time")
	}
}
