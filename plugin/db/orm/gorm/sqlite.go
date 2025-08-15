package gorm

import (
	"database/sql"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewSqliteDB creates a new GORM DB instance using a database connection.
func NewSqliteDB(db *sql.DB, opts ...gorm.Option) (*gorm.DB, error) {
	return gorm.Open(sqlite.New(sqlite.Config{
		Conn: db,
	}), opts...)
}

// NewSqliteDBByDSN creates a new GORM DB instance using a database connection string.
func NewSqliteDBByDSN(dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), opts...)
}

// NewSqliteDBByConfig creates a new GORM DB instance using a database configuration.
func NewSqliteDBByConfig(dbConfig sqlite.Config, opts ...gorm.Option) (*gorm.DB, error) {
	return gorm.Open(sqlite.New(dbConfig), opts...)
}
