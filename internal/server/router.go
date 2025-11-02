package server

import (
	"context"
	"log"
	"net/http"

	"github.com/TonyGLL/gofetch/internal/analysis"
	"github.com/TonyGLL/gofetch/internal/search"
	"github.com/TonyGLL/gofetch/internal/server/handler"
	"github.com/TonyGLL/gofetch/internal/server/middleware"
	"github.com/TonyGLL/gofetch/internal/storage"
)

func NewRouter() *http.ServeMux {
	// --- Dependency Injection ---

	// 1. Create and connect to the database store.
	store := storage.NewMongoStore()
	if err := store.Connect(context.Background()); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// 2. Create the analyzer.
	analyzer := analysis.NewEnglishAnalyzer()

	// 3. Create the searcher with its dependencies.
	searcher := search.NewSearcher(analyzer, store)

	// 4. Create the search handler with its dependency.
	searchHandler := &handler.Search{
		Searcher: searcher,
	}

	// --- Routing ---

	mux := http.NewServeMux()

	v1 := http.NewServeMux()
	v1.Handle("GET /search", searchHandler)

	// Chain middleware
	v1WithMiddleware := middleware.Chain(
		v1,
		middleware.CORS, // Add CORS middleware
		middleware.Logging,
	)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1WithMiddleware))
	// --- Static File Server for UI ---
	fs := http.FileServer(http.Dir("./ui"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./ui/index.html")
	})

	return mux
}
