package main

import (
	"SERVER/API"
	"fmt"
	"net/http"
)

func main() {
	cache := API.NewLRUCache(2)

	http.HandleFunc("/get", cache.GET)
	http.HandleFunc("/post", cache.POST)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
