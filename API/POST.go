package API

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (cache *LRUCache) POST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	key := r.FormValue("key")
	value := r.FormValue("value")
	expirationString := r.FormValue("expiration")
	var expiration time.Time
	if expirationString != "" {
		var err error
		seconds, err := strconv.Atoi(expirationString)
		if err != nil {
			http.Error(w, "expiration should be in integer", http.StatusBadRequest)
			return
		}

		expiration = time.Now().Add(time.Second * time.Duration(seconds))
	}

	if key == "" || value == "" {
		http.Error(w, "Both key and value parameters are required", http.StatusBadRequest)
		return
	}

	resp := cache.Set(key, value, expiration)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("error while decoding the response, err : " + err.Error())
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(jsonResp)
}

// Set sets a key-value pair in the cache with an expiration time.
func (c *LRUCache) Set(key, value string, expiration time.Time) Response {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// If key already exists, update the value and expiration time, move it to the front.
	if entry, exists := c.cache[key]; exists {
		entry.value = value
		entry.expiration = expiration
		c.lruList.MoveToFront(entry.listElement)
		return Response{Key: key, Value: value, Expiration: expiration}
	}

	// If capacity is reached, remove the least recently used item from cache.
	if len(c.cache) >= c.capacity {
		// Get the least recently used key.
		lastNode := c.lruList.Back()
		k := lastNode.Value.(string)
		// Remove it from cache and the list.
		delete(c.cache, k)

		c.lruList.Remove(lastNode)
	}

	// Add the new key-value pair to the cache with expiration time.
	listElement := c.lruList.PushFront(key)
	c.cache[key] = CacheEntry{value: value, expiration: expiration, listElement: listElement}

	return Response{Key: key, Value: value, Expiration: expiration}
}