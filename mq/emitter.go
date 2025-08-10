package mq

import (
	"context"
	"encoding/json"
	"log"
	"naevis/globals"
	"naevis/models"
	"naevis/search"
)

// Index represents an indexing-related message to be emitted.
type Index struct {
	EntityType string `json:"entity_type"`
	Method     string `json:"method"`
	EntityId   string `json:"entity_id"`
	ItemId     string `json:"item_id"`
	ItemType   string `json:"item_type"`
}

func StartIndexingWorker() {
	ctx := context.Background()
	sub := globals.RedisClient.Subscribe(ctx, "indexing-events")
	ch := sub.Channel()

	log.Println("[IndexingWorker] Listening for indexing events...")

	for msg := range ch {
		var event models.Index
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Printf("[IndexingWorker] Failed to parse event: %v", err)
			continue
		}
		log.Printf("[IndexingWorker] Processing event=%+v", event)

		if err := search.IndexDatainRedis(ctx, event); err != nil {
			log.Printf("[IndexingWorker] IndexDatainRedis error: %v", err)
		} else {
			log.Println("[IndexingWorker] Indexing complete")
		}
	}
}

// Printer sends raw JSON to an external endpoint (placeholder).
// Replace with your actual QUIC or HTTP implementation when ready.
/*
func Printer(jsonData []byte) error {
	start := time.Now()
	resp, err := http.Post(SERP_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Server Response: %s\n", string(body))
	fmt.Printf("Execution Time: %v\n", elapsed)

	return nil
}
*/
