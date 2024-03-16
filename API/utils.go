package API

import (
	"fmt"
	"time"
)

// cleanupExpiredEntries periodically checks for and removes expired cache entries.
func (c *LRUCache) cleanupExpiredEntries() {
	cleanupTicker := time.NewTicker(time.Minute) // Adjust the interval as needed
	defer cleanupTicker.Stop()

	for {
		<-cleanupTicker.C
		c.mutex.Lock()
		for key, entry := range c.cache {
			if entry.expired() {
				delete(c.cache, key)
				c.lruList.Remove(entry.listElement)
				fmt.Println("key : " + key + " removed at " + time.Now().Format(time.DateTime))
			}
		}
		c.mutex.Unlock()
	}
}

// expired checks if the cache entry has expired.
func (e CacheEntry) expired() bool {
	return !e.expiration.IsZero() && time.Now().After(e.expiration)
}
