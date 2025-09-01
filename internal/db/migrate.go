package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func Migrate(ctx context.Context, sqlDB *sql.DB, migrationsDir string) error {
	_, err := sqlDB.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version VARCHAR(255) PRIMARY KEY,
            applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
    `)
	if err != nil {
		return err
	}

	applied := map[string]bool{}
	rows, err := sqlDB.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return err
		}
		applied[v] = true
	}

	dir := os.DirFS(migrationsDir)
	return fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}
		if applied[path] {
			return nil
		}

		sqlBytes, readErr := fs.ReadFile(dir, path)
		if readErr != nil {
			return readErr
		}

		if _, execErr := sqlDB.ExecContext(ctx, string(sqlBytes)); execErr != nil {
			return fmt.Errorf("migration %s failed: %w", path, execErr)
		}

		if _, insErr := sqlDB.ExecContext(ctx,
			`INSERT INTO schema_migrations(version) VALUES(?)`, path); insErr != nil {
			return insErr
		}
		return nil
	})
}

var ErrNoMigrations = errors.New("no migrations directory")
