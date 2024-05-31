package main

import (
	"context"
	db "go-db-balancer/internal/adapters"
	"go-db-balancer/internal/pkg/app"
	"log/slog"
	"os"
)

func main(){

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbAdapter, err := db.NewAdapter(*logger)

	if err != nil {
		logger.Error("Unable to Reach DynamoDB. Error: %v", err)
	}

	app := app.NewApp(dbAdapter, *logger)

	app.StartDatabaseBalancing(context.TODO())
}