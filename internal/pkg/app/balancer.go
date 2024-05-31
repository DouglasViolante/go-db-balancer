package app

import (
	"context"
	"go-db-balancer/internal/pkg/domain"
	utils "go-db-balancer/internal/pkg/utils"
	"go-db-balancer/internal/ports"
	"log/slog"
	"math"
)

type App struct {
	db 	ports.DatabasePort
	logger slog.Logger
}

func NewApp(db ports.DatabasePort, logger slog.Logger) *App {
	return &App{
		db: 	db,
		logger: logger,
	}
}

func (app App) VerifyDatabases(ctx context.Context, tableNames []string) ([]domain.Table) {
	tablesInfo, _ := app.db.GetTablesItemCount(ctx, tableNames)

	if tablesInfo[0].ItemCount == tablesInfo[1].ItemCount {
		app.logger.Error("DynamoDB Tables %s and %s Already Balanced! Aborting...", tableNames[0], tableNames[1])
	}
	return tablesInfo
}

func (app App) StartDatabaseBalancing(ctx context.Context) error {

	tableNames := []string{"db_A", "db_B"}

	tablesInfo := app.VerifyDatabases(ctx, tableNames)

	biggestTablePos 	:= utils.FindMaxPos(tablesInfo)
	smallestTablePos 	:= int(math.Abs(float64(biggestTablePos - 1)))

	balancingThreshold := math.Round(float64(tablesInfo[biggestTablePos].ItemCount) * 0.7)

	if tablesInfo[smallestTablePos].ItemCount >= int(balancingThreshold) {
		app.logger.Info("DynamoDB Tables %s and %s Balancing is NOT Needed, Difference is still within Threshold! Aborting...", tableNames[0], tableNames[1])
		return nil
	}

	itemCountTarget := (tablesInfo[0].ItemCount + tablesInfo[1].ItemCount) / 2
	app.logger.Info("The Item Count Target for Balancing is %d", itemCountTarget)

	numItemsToGet := int(math.Abs(float64(itemCountTarget - tablesInfo[biggestTablePos].ItemCount)))
	app.logger.Info("The Total Items to be Retrieved from %s, is %d", tablesInfo[biggestTablePos].TableName, numItemsToGet)

	retrievedRecords, _ := app.db.GetItems(ctx, numItemsToGet, tablesInfo[biggestTablePos].TableName)

	app.db.PutItems(ctx, retrievedRecords, tablesInfo[smallestTablePos].TableName)
	app.db.DeleteItems(ctx, retrievedRecords, tablesInfo[biggestTablePos].TableName)

	app.logger.Info("Balacing of Tables %s and %s has been Completed!", tableNames[0], tableNames[1])

	return nil
}