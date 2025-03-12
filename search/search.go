package search

import (
	"log"
	"naevis/handlers"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Search handler (fetches based on active tab)
func SearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	entityType := ps.ByName("entityType") // Extract active tab type
	log.Println("Received search request for:", entityType)

	query := r.URL.Query().Get("query")

	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}
	handlers.GetResultsByTypeHandler(w, r, entityType, query)
}
