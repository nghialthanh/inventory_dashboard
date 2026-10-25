package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/*.sql
var migrationSQL embed.FS

// Migrate applies the schema and the catalog seed. Statements are idempotent.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := fs.ReadDir(migrationSQL, "sql")
	if err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	for _, entry := range entries {
		raw, err := migrationSQL.ReadFile("sql/" + entry.Name())
		if err != nil {
			return err
		}
		results, err := conn.Conn().PgConn().Exec(ctx, string(raw)).ReadAll()
		if err != nil {
			return fmt.Errorf("%s: %w", entry.Name(), err)
		}
		for _, result := range results {
			if err := result.Err; err != nil {
				return fmt.Errorf("%s: %w", entry.Name(), err)
			}
		}
	}
	return nil
}
