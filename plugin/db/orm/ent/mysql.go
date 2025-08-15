// Package ent provides a GORM-like interface for interacting with databases using Ent.
package ent

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"contrib.go.opencensus.io/integrations/ocsql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/go-sql-driver/mysql"
)

// NewMysqlDB creates a new MySQL database connection using a database connection.
func NewMysqlDB(db *sql.DB) *entsql.Driver {
	return entsql.OpenDB("mysql", db)
}

type connector struct {
	dsn string
}

func (c connector) Connect(context.Context) (driver.Conn, error) {
	return c.Driver().Open(c.dsn)
}

func (connector) Driver() driver.Driver {
	return ocsql.Wrap(
		mysql.MySQLDriver{},
		ocsql.WithAllTraceOptions(),
		ocsql.WithRowsClose(false),
		ocsql.WithRowsNext(false),
		ocsql.WithDisableErrSkip(true),
	)
}

// NewMysqlDBByDSN creates a new MySQL database connection using a DSN string.
func NewMysqlDBByDSN(dsn string) *entsql.Driver {
	db := sql.OpenDB(connector{dsn})
	return entsql.OpenDB(dialect.MySQL, db)
}
