package cli

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
)

const sqlListUpdateQuery = `
      UPDATE client SET name=$2 WHERE id=$1
`

func (repo *ClientRepository) Update(ctx context.Context, dto repository.ClientDTO) (pgtype.UUID, error) {
	connection := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	_, err := connection.Exec(ctx, sqlListUpdateQuery, dto.ID, dto.Name)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return dto.ID, nil
}
