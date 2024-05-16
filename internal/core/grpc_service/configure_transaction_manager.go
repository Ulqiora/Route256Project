package grpc_service

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/config"
	"github.com/Ulqiora/Route256Project/internal/database/postgresql"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

func ConfigureTransactionManager(ctx context.Context, config *config.Config) (*manager.Manager, *postgresql.Database, error) {
	database, err := postgresql.NewDb(ctx, config.PostgresDsn)
	if err != nil {
		return nil, nil, err
	}
	trManager := manager.Must(trmpgx.NewDefaultFactory(database.GetPool(ctx)))
	return trManager, database, nil
}
