package ports

import (
	"context"
	"go-db-balancer/internal/pkg/domain"
)

type DatabasePort interface {
	GetTablesItemCount(ctx context.Context, tableName []string) ([]domain.Table, error)
	GetItems(ctx context.Context, totalItems int, sourceTableName string) ([]domain.Specs, error)
	PutItems(ctx context.Context, recordsToMove []domain.Specs, putTargetTableName string) error
	DeleteItems(ctx context.Context, recordsToDelete []domain.Specs, deleteTargetTableName string) error
}