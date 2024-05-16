package cli

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/repository"
)

const sqlListClientQuery string = `
     SELECT
         id,name
     FROM client
     WHERE deleted=false
`

func (repo *ClientRepository) List(ctx context.Context) ([]repository.ClientDTO, error) {
	connection := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	rows, err := connection.Query(ctx, sqlListClientQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []repository.ClientDTO
	for rows.Next() {
		var client repository.ClientDTO
		err := client.LoadFromRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, client)
	}
	if repo.cache != nil {
		for _, client := range result {
			err = repo.cache.Set(ctx, string(client.ID.Bytes[:]), client)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}
