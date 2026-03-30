package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

func OpenDatabase(dbType, dsn string) (*sql.DB, Dialect, error) {
	var dialect Dialect
	switch dbType {
	case "sqlite":
		dialect = SQLiteDialect{}
		if dsn == "" || dsn == ":memory:" {
			dsn = "file::memory:?cache=shared&_pragma=foreign_keys(1)"
		} else {
			// File-based SQLite: add WAL mode and foreign keys
			if !strings.Contains(dsn, "_pragma") {
				dsn = "file:" + dsn + "?_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)"
			}
		}
	case "postgres":
		dialect = PostgresDialect{}
		if dsn == "" {
			return nil, nil, fmt.Errorf("HABITCLAW_DB_DSN is required for postgres")
		}
	case "mysql":
		dialect = MySQLDialect{}
		if dsn == "" {
			return nil, nil, fmt.Errorf("HABITCLAW_DB_DSN is required for mysql")
		}
		if !strings.Contains(dsn, "parseTime") {
			if strings.Contains(dsn, "?") {
				dsn += "&parseTime=true"
			} else {
				dsn += "?parseTime=true"
			}
		}
	default:
		return nil, nil, fmt.Errorf("unsupported database type: %s (supported: sqlite, postgres, mysql)", dbType)
	}

	db, err := sql.Open(dialect.DriverName(), dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, dialect, nil
}
