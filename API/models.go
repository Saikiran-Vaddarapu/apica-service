package API

import (
	"container/list"
	"sync"
	"time"
)

type Response struct {
	Key        string    `json:"key"`
	Value      string    `json:"value"`
	Expiration time.Time `json:"expiration"`
}

// CacheEntry represents a key-value pair with expiration time.
type CacheEntry struct {
	value       string
	expiration  time.Time
	listElement *list.Element
}

// LRUCache represents a least recently used cache with expiration time.
type LRUCache struct {
	capacity int
	cache    map[string]CacheEntry
	lruList  *list.List
	mutex    sync.Mutex
}
