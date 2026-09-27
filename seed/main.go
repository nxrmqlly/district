package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nxrmqlly/district/store"
)

func main() {
	godotenv.Load()

	ctx := context.Background()

	poolCfg, err := pgxpool.ParseConfig(os.Getenv("POSTGRES_CONNSTR"))
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	q := store.New(pool)

	u, err := q.NewUser(ctx, store.NewUserParams{
		Username:   "superadmin",
		Email:      "what@lol.lol",
		PasswdHash: "67890897865432456789897675764653",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(u.ID.String())
}
