package server

import (
	"net/http"

	"github.com/TonyGLL/gofetch/internal/server/handler"
	"github.com/TonyGLL/gofetch/internal/server/middleware"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	v1 := http.NewServeMux()
	v1.HandleFunc("GET /search", handler.SearchHandler)

	v1WithMiddleware := middleware.Chain(
		v1,
		middleware.Logging,
	)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1WithMiddleware))
	return mux
}
