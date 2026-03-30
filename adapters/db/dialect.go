package db

import (
	"strconv"
	"strings"
)

// Dialect abstracts database-specific SQL differences.
type Dialect interface {
	DriverName() string
	Rebind(query string) string
	MigrationsDir() string
}

type SQLiteDialect struct{}

func (d SQLiteDialect) DriverName() string         { return "sqlite" }
func (d SQLiteDialect) Rebind(query string) string { return query }
func (d SQLiteDialect) MigrationsDir() string      { return "sqlite" }

type PostgresDialect struct{}

func (d PostgresDialect) DriverName() string { return "pgx" }
func (d PostgresDialect) Rebind(query string) string {
	var buf strings.Builder
	n := 1
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			buf.WriteString("$" + strconv.Itoa(n))
			n++
		} else {
			buf.WriteByte(query[i])
		}
	}
	return buf.String()
}
func (d PostgresDialect) MigrationsDir() string { return "postgres" }

type MySQLDialect struct{}

func (d MySQLDialect) DriverName() string         { return "mysql" }
func (d MySQLDialect) Rebind(query string) string { return query }
func (d MySQLDialect) MigrationsDir() string      { return "mysql" }
