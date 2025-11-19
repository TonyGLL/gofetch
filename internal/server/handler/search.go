package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TonyGLL/gofetch/internal/search"
	"github.com/TonyGLL/gofetch/pkg/storage"
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

	page := r.URL.Query().Get("page")
	if page == "" {
		http.Error(w, "query parameter 'page' is missing", http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		http.Error(w, "query parameter 'limit' is missing", http.StatusBadRequest)
		return
	}

	// 2. Perform the search using the injected searcher.
	pageInt64, err := strconv.ParseInt(page, 0, 0)
	if err != nil {
		// Log the error internally
		// In a real app, you'd use a structured logger.
		// log.Printf("error during search: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	limitInt64, err := strconv.ParseInt(limit, 0, 0)
	if err != nil {
		// Log the error internally
		// In a real app, you'd use a structured logger.
		// log.Printf("error during search: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	results, err := s.Searcher.Search(r.Context(), query, storage.GetDocumentsFilter{
		Page:  pageInt64,
		Limit: limitInt64,
	})
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
		log.Printf("error encoding response: %v", err)
	}
}
