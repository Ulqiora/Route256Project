package cli

import (
	"context"
	"fmt"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
)

const sqlCreateClientQuery string = `
	INSERT INTO client(name) VALUES ($1) RETURNING id
`

func (repo *ClientRepository) Create(ctx context.Context, client repository.ClientDTO) (pgtype.UUID, error) {
	connection := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	//id
	err := connection.QueryRow(ctx, sqlCreateClientQuery, client.Name).Scan(&client.ID)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("%s: %w", repository.ErrorObjectNotCreated.Error(), err)
	}
	return client.ID, nil
}
