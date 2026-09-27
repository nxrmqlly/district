package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nxrmqlly/district/app"
	"github.com/nxrmqlly/district/store"
)

func main() {
	godotenv.Load()

	ctx := context.Background()

	store.Migrate(ctx, os.Getenv("POSTGRES_CONNSTR"))

	fmt.Println("hello world")

	srv := http.Server{
		Addr:    os.Getenv("BIND_ADDR"),
		Handler: app.New(),
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
