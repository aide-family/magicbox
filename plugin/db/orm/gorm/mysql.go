package gorm

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMysqlDB creates a new GORM DB instance using a database connection.
func NewMysqlDB(db *sql.DB, opts ...gorm.Option) (*gorm.DB, error) {
	return gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), opts...)
}

// NewMysqlDBByDSN creates a new GORM DB instance using a DSN string.
func NewMysqlDBByDSN(dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), opts...)
}

// NewMysqlDBByConfig creates a new GORM DB instance using a database configuration.
func NewMysqlDBByConfig(dbConfig mysql.Config, opts ...gorm.Option) (*gorm.DB, error) {
	return gorm.Open(mysql.New(dbConfig), opts...)
}
