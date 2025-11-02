package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TonyGLL/gofetch/internal/search"
)

// Search is the handler for the search endpoint.
// It holds a dependency to the Searcher interface.
type Search struct {
	Searcher search.Searcher
}

// ServeHTTP handles the HTTP request for a search.
func (s *Search) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Get the query from the URL parameters.
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is missing", http.StatusBadRequest)
		return
	}

	// 2. Perform the search using the injected searcher.
	results, err := s.Searcher.Search(r.Context(), query)
	if err != nil {
		// Log the error internally
		// In a real app, you'd use a structured logger.
		// log.Printf("error during search: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// 3. Write the JSON response.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		// This error is harder to handle as the headers might already be written.
		// log.Printf("error encoding response: %v", err)
	}
}
