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

type CacheEntry struct {
	value       string
	expiration  time.Time
	listElement *list.Element
}

type LRUCache struct {
	capacity int
	cache    map[string]CacheEntry
	lruList  *list.List
	mutex    sync.Mutex
}
