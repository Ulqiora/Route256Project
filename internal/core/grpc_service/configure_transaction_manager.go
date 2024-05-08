package grpc_service

import (
	"context"

	"homework/internal/config"
	"homework/internal/database/postgresql"
	"homework/internal/database/transaction_manager"
)

func ConfigureTransactionManager(ctx context.Context, config *config.Config) (*transaction_manager.TransactionManager, error) {
	database, err := postgresql.NewDb(ctx, config.PostgresDsn)
	if err != nil {
		return nil, err
	}
	manager := transaction_manager.New(database)
	return manager, nil
}
