// Package sqlite provides a SQLite driver.
package sqlite

import (
	"database/sql"

	"github.com/aide-family/magicbox/plugin/db"
)

var _ db.Driver = (*initializer)(nil)

// NewDBDriver creates a new SQLite driver.
func NewDBDriver(dsn string) db.Driver {
	return &initializer{dsn: dsn}
}

type initializer struct {
	dsn string
}

// New creates a new SQLite database connection.
func (i *initializer) New() (*sql.DB, error) {
	return sql.Open("sqlite3", i.dsn)
}
