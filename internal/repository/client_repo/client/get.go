package cli

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
)

const sqlGetClientQuery = `
    SELECT
        id,name
    FROM client
    WHERE id = $1 AND client.deleted=false
`

func (repo *ClientRepository) GetByID(ctx context.Context, id pgtype.UUID) (repository.ClientDTO, error) {
	var client repository.ClientDTO
	if repo.cache != nil {
		/*idBytes := id.Get().([16]byte)
		  bytes, err := repo.cache.Get(ctx, string(idBytes[:]))
		  if err == nil {
		      return client, err
		  }
		  err = json.Unmarshal(bytes, &client)
		  if err == nil {
		      return client, nil
		  }
		  slog.Info(err.Error())

		*/
	}

	connection := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	row := connection.QueryRow(ctx, sqlGetClientQuery, id)
	if client.LoadFromRow(row) != nil {
		return client, repository.ErrorObjectNotCreated
	}
	return client, nil
}
