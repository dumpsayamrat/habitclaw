package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"
)

//go:embed migrations/sqlite/*.sql migrations/postgres/*.sql migrations/mysql/*.sql
var migrationFS embed.FS

func Migrate(database *sql.DB, dialect Dialect) error {
	createTable := `CREATE TABLE IF NOT EXISTS schema_migrations (
        version    TEXT PRIMARY KEY,
        applied_at TEXT NOT NULL
    )`
	if _, err := database.Exec(createTable); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	dir := "migrations/" + dialect.MigrationsDir()
	entries, err := fs.ReadDir(migrationFS, dir)
	if err != nil {
		return fmt.Errorf("read migrations dir %s: %w", dir, err)
	}

	var filenames []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			filenames = append(filenames, e.Name())
		}
	}
	sort.Strings(filenames)

	applied := make(map[string]bool)
	rows, err := database.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("query applied migrations: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return err
		}
		applied[v] = true
	}

	for _, filename := range filenames {
		if applied[filename] {
			continue
		}
		content, err := fs.ReadFile(migrationFS, dir+"/"+filename)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", filename, err)
		}

		tx, err := database.Begin()
		if err != nil {
			return fmt.Errorf("begin tx for %s: %w", filename, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("execute migration %s: %w", filename, err)
		}

		if _, err := tx.Exec(dialect.Rebind("INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)"), filename, time.Now().Format(time.RFC3339)); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %s: %w", filename, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", filename, err)
		}
	}

	return nil
}
