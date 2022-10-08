package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"hack/internal/app"
	"hack/internal/app/clientRouter"
	"hack/internal/app/parserRouter"
	"hack/internal/pkg/parsers/consultantParser"
	"hack/internal/pkg/parsers/lentaParser"
	store2 "hack/internal/pkg/store"
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
	Prouter := parserRouter.NewRouter(lParser, cParser)
	CRouter := clientRouter.NewRouter(store)
	router := app.NewRouter(Prouter, CRouter)
	router.SetUpRouter()
	router.Run()
}
