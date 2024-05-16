package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Ulqiora/Route256Project/internal/config"
	"github.com/Ulqiora/Route256Project/internal/database/postgresql"
)

type TestDatabase struct {
	postgresql.PGXDatabase
}

func (r *TestDatabase) TruncateAll(ctx context.Context) error {
	var tables []string
	err := r.Select(ctx, &tables, "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name != 'goose_db_version'")
	if err != nil {
		return err
	}
	if len(tables) == 0 {
		return errors.New("run migration first")
	}
	q := fmt.Sprintf("TRUNCATE table %s", strings.Join(quoteEach(tables), ","))
	if _, err := r.Exec(ctx, q); err != nil {
		return err
	}
	return nil
}
func (r *TestDatabase) TruncateTable(ctx context.Context, tablename string) error {
	q := fmt.Sprintf("TRUNCATE table %s CASCADE", tablename)
	if _, err := r.Exec(ctx, q); err != nil {
		return err
	}
	return nil
}
func quoteEach(strs []string) []string {
	quoted := make([]string, len(strs))
	for i, s := range strs {
		quoted[i] = fmt.Sprintf("'%s'", s)
	}
	return quoted
}

func NewTestDatabase(ctx context.Context, config *config.Config) (*TestDatabase, error) {
	db, err := postgresql.NewDb(ctx, config.PostgresDsn)
	if err != nil {
		return nil, err
	}
	return &TestDatabase{db}, nil
}
