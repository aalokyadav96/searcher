package search

import (
	"encoding/json"
	"fmt"
	"log"
	"naevis/initdb"
	"strings"
	"time"
)

// -------------------------
// REDIS CONTEXT
// -------------------------

// var ctx = context.Background()

// -------------------------
// TRIE FOR AUTOCOMPLETE
// -------------------------

// // TrieNode represents a node in the Trie.
// type TrieNode struct {
// 	Children map[rune]*TrieNode
// 	IsWord   bool
// }

// // NewTrieNode creates a new TrieNode.
// func NewTrieNode() *TrieNode {
// 	return &TrieNode{Children: make(map[rune]*TrieNode)}
// }

// AddWord inserts a word into the Trie and saves it in Redis.
func (t *TrieNode) AddWord(word string) error {
	node := t
	for _, ch := range word {
		if node.Children == nil {
			node.Children = make(map[rune]*TrieNode)
		}
		if _, exists := node.Children[ch]; !exists {
			node.Children[ch] = NewTrieNode()
		}
		node = node.Children[ch]
	}
	node.IsWord = true

	// Also store the word in Redis for persistence.
	return SaveAutocompleteWord(word)
}

// SaveAutocompleteWord stores an autocomplete word in Redis.
func SaveAutocompleteWord(word string) error {
	key := fmt.Sprintf("autocomplete:%s", word)
	return initdb.RedisClient.Set(ctx, key, 1, 0).Err() // No expiration, acts as a unique set.
}

// GetWordsWithPrefix fetches autocomplete suggestions from Redis.
func GetWordsWithPrefix(prefix string) ([]string, error) {
	prefix = strings.ToLower(prefix)
	cacheKey := fmt.Sprintf("autocomplete_cache:%s", prefix)

	// Check if results exist in cache.
	if cached, err := GetCachedAutocompleteResults(cacheKey); err == nil {
		return cached, nil
	}

	// Use Redis SCAN command to find matching keys.
	iter := initdb.RedisClient.Scan(ctx, 0, fmt.Sprintf("autocomplete:%s*", prefix), 100).Iterator()
	var results []string
	for iter.Next(ctx) {
		results = append(results, strings.TrimPrefix(iter.Val(), "autocomplete:"))
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	// Cache the results in Redis for faster future lookups.
	if err := CacheAutocompleteResults(cacheKey, results); err != nil {
		log.Println("Error caching autocomplete results:", err)
	}
	return results, nil
}

// -------------------------
// REDIS CACHING HELPERS
// -------------------------

// CacheAutocompleteResults stores autocomplete suggestions in Redis with an expiration.
func CacheAutocompleteResults(key string, results []string) error {
	data, err := json.Marshal(results)
	if err != nil {
		return err
	}
	return initdb.RedisClient.Set(ctx, key, data, time.Hour).Err()
}

// GetCachedAutocompleteResults retrieves cached autocomplete suggestions.
func GetCachedAutocompleteResults(key string) ([]string, error) {
	data, err := initdb.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var results []string
	if err = json.Unmarshal([]byte(data), &results); err != nil {
		return nil, err
	}
	return results, nil
}

// -------------------------
// AUTOCOMPLETE HANDLER
// -------------------------

// // Autocompleter handles HTTP requests for autocomplete suggestions.
// func Autocompleter(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	prefix := r.URL.Query().Get("prefix")
// 	if prefix == "" {
// 		http.Error(w, "Search prefix is required", http.StatusBadRequest)
// 		return
// 	}
// 	prefix = strings.ToLower(prefix)

// 	// Retrieve suggestions.
// 	results, err := GetWordsWithPrefix(prefix)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Error retrieving autocomplete: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Convert to JSON and respond.
// 	response, err := json.Marshal(results)
// 	if err != nil {
// 		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(response)
// }
