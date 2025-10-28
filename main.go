package main

import (
	"fmt"

	"github.com/TonyGLL/gofetch/cmd/indexer"
)

func main() {
	indexer.Execute()
	fmt.Println("Indexer has finished execution.")
}
