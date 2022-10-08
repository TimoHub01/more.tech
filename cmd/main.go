package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"hack/internal/app/pkg/parsers"
	"hack/internal/app/pkg/parsers/consultantParser"
	"hack/internal/app/pkg/parsers/lentaParser"
	store2 "hack/internal/app/pkg/store"
	"os"
)

func main() {
	databaseUrl := "postgres://admin:more-tech@localhost:15432/postgres"
	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Print(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	store := store2.NewPostgresStore(dbPool)
	lParser := lentaParser.NewParser(store)
	cParser := consultantParser.NewParser(store)
	router := parsers.NewRouter(lParser, cParser)
	router.SetUpRouter()
	router.Run()
}
