package search

import (
	"encoding/json"
	"log"
	"net/http"

	"naevis/initdb"
	"naevis/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = initdb.CTX

// GetResultsByTypeHandler handles requests to /{ENTITY_TYPE}?query=QUERY
func GetResultsByTypeHandler(w http.ResponseWriter, r *http.Request, entityType string, query string) {
	log.Println("Entity Type:", entityType)

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get results based on the reverse index
	results := GetResultsOfType(entityType, query)

	// Convert results slice (or map for "all") to JSON and send response
	response, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetResultsOfType fetches results based on entity type using a reverse index.
// Instead of mapping to a generic result, we now decode directly into the Event or Place structs.
func GetResultsOfType(entityType string, query string) interface{} {
	switch entityType {
	case "events":
		eventIDs := GetIndexResults("events", query)
		events := []structs.Event{}
		for _, id := range eventIDs {
			filter := bson.M{"eventid": id}
			var ev structs.Event
			err := FetchAndDecode("events", filter, &ev)
			if err != nil {
				log.Println("Error fetching event:", err)
				continue
			}
			events = append(events, ev)
		}
		return events

	case "places":
		placeIDs := GetIndexResults("places", query)
		places := []structs.Place{}
		for _, id := range placeIDs {
			filter := bson.M{"placeid": id}
			var p structs.Place
			err := FetchAndDecode("places", filter, &p)
			if err != nil {
				log.Println("Error fetching place:", err)
				continue
			}
			places = append(places, p)
		}
		return places

	case "all":
		// For the "all" case, you might want to return both events and places.
		eventIDs := GetIndexResults("events", query)
		events := []structs.Event{}
		for _, id := range eventIDs {
			filter := bson.M{"eventid": id}
			var ev structs.Event
			err := FetchAndDecode("events", filter, &ev)
			if err != nil {
				log.Println("Error fetching event:", err)
				continue
			}
			events = append(events, ev)
		}

		placeIDs := GetIndexResults("places", query)
		places := []structs.Place{}
		for _, id := range placeIDs {
			filter := bson.M{"placeid": id}
			var p structs.Place
			err := FetchAndDecode("places", filter, &p)
			if err != nil {
				log.Println("Error fetching place:", err)
				continue
			}
			places = append(places, p)
		}
		return map[string]interface{}{
			"events": events,
			"places": places,
		}

	default:
		return nil
	}
}

// GetIndexResults simulates a reverse index retriever.
// In your real implementation, this function would query your reverse index.
func GetIndexResults(entityType string, query string) []string {
	// For demonstration purposes, we return hard-coded IDs.
	// You can modify this to perform your reverse-index lookup.
	switch entityType {
	case "events":
		return []string{"Bz7TnUhPbr0_4g"}
	case "places":
		return []string{"JQLCD3D_t0mhR_"}
	default:
		return []string{}
	}
}

// FetchAndDecode retrieves a document from MongoDB and decodes it into the provided output struct.
func FetchAndDecode(collectionName string, filter bson.M, out interface{}) error {
	collection := initdb.Client.Database("eventdb").Collection(collectionName)
	projection, exists := Projections[collectionName]
	if !exists {
		projection = bson.M{}
	}
	opts := options.FindOne().SetProjection(projection)
	return collection.FindOne(ctx, filter, opts).Decode(out)
}
