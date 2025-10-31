package server

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	srv *http.Server
}

func NewServer(handler http.Handler, port string) *Server {
	return &Server{
		srv: &http.Server{
			Addr:           fmt.Sprintf(":%s", port),
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

func (s *Server) Start() error {
	fmt.Printf("Starting server on %s\n", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown() error {
	fmt.Println("Shutting down server...")
	return s.srv.Close()
}
