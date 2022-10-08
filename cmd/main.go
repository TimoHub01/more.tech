package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"hack/internal/app"
	"hack/internal/app/clientRouter"
	"hack/internal/app/parserRouter"
	"hack/internal/pkg/parsers/consultantParser"
	"hack/internal/pkg/parsers/kommersantParser"
	"hack/internal/pkg/parsers/lentaParser"
	"hack/internal/pkg/parsers/newslab"
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
	KParser := kommersantParser.NewParser(store)
	NParser := newslab.NewParser(store)
	ParseRouter := parserRouter.NewRouter(NParser, lParser, cParser, KParser)

	CPRouter := clientRouter.NewRouter(store)
	router := app.NewRouter(ParseRouter, CPRouter)
	router.SetUpRouter()
	router.Run()
}
