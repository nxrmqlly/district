package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nxrmqlly/district/app"
	"github.com/nxrmqlly/district/store"
)

func main() {
	godotenv.Load()

	ctx := context.Background()

	store.Migrate(ctx, os.Getenv("POSTGRES_CONNSTR"))

	poolCfg, err := pgxpool.ParseConfig(os.Getenv("POSTGRES_CONNSTR"))
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	router, err := app.NewRouter(store.New(pool))
	if err != nil {
		log.Fatal(err)
	}

	srv := http.Server{
		Addr:    os.Getenv("BIND_ADDR"),
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
