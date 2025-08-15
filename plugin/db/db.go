// Package db provides a database connection pool.
package db

import "database/sql"

// New creates a new database connection pool by driver.
func New(driver Driver) (*sql.DB, error) {
	return driver.New()
}

// Driver is the interface that wraps the New method.
type Driver interface {
	New() (*sql.DB, error)
}
