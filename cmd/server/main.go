package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/TonyGLL/gofetch/internal/server"
)

func main() {
	handler := server.NewRouter()
	srv := server.NewServer(handler, "8080")

	go func() { log.Fatal(srv.Start()) }()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	if err := srv.Shutdown(); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
