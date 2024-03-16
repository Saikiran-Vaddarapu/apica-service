package API

import (
	"container/list"
	"encoding/json"
	"fmt"
	"net/http"
)

func NewLRUCache(capacity int) *LRUCache {
	cache := &LRUCache{
		capacity: capacity,
		cache:    make(map[string]CacheEntry),
		lruList:  list.New(),
	}
	go cache.cleanupExpiredEntries()
	return cache
}

func (cache *LRUCache) GET(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key parameter is required", http.StatusBadRequest)
		return
	}

	resp, exists := cache.Retrieve(key)
	if !exists {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("error while decoding the response, err : " + err.Error())
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResp)
}

func (c *LRUCache) Retrieve(key string) (Response, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.cache[key]
	if !exists || entry.expired() {
		return Response{}, false
	}

	fmt.Println(entry, exists)

	c.lruList.MoveToFront(entry.listElement)

	return Response{Key: key, Value: entry.value, Expiration: entry.expiration}, true
}
