package indexer

import (
	"context"
	"flag"
	"fmt"

	"github.com/TonyGLL/gofetch/internal/storage"
)

func Execute() {
	store := storage.NewMongoStore()
	ctx := context.Background()
	if err := store.Connect(ctx); err != nil {
		panic(err)
	}
	defer func() {
		if err := store.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	args := flag.String("path", "data/index", "Path to store the index data")
	flag.Parse()
	fmt.Println("args:", *args)
}
