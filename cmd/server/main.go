package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Server is running on port 8080")
	log.Fatal(s.ListenAndServe())
}
