// Package mysql provides a MySQL driver.
package mysql

import (
	"database/sql"

	"github.com/aide-family/magicbox/plugin/db"
)

var _ db.Driver = (*initializer)(nil)

func NewDBDriver(dsn string) db.Driver {
	return &initializer{dsn: dsn}
}

type initializer struct {
	dsn string
}

func (i *initializer) New() (*sql.DB, error) {
	return sql.Open("mysql", i.dsn)
}
