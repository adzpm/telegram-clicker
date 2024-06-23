package main

import (
	"context"
	"github.com/adzpm/tg-clicker/internal/storage"

	zap "go.uber.org/zap"

	api "github.com/adzpm/tg-clicker/internal/api"
)

func main() {
	ctx, can := context.WithCancel(context.Background())
	defer can()

	lgr, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer lgr.Sync()

	str, err := storage.NewStorage(lgr, &storage.Config{
		Path: "tgc.db",
	})

	if err != nil {
		panic(err)
	}

	a := api.NewAPI(lgr, str, &api.Config{
		Host:    "127.0.0.1",
		Port:    "8080",
		WebPath: "./web/",
	})

	if err = a.Start(ctx); err != nil {
		panic(err)
	}
}
