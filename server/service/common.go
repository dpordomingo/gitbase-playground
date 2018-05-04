package service

import "database/sql"

// SQLDB mocks a sql.DB
type SQLDB interface {
	Close() error
	Query(query string, args ...interface{}) (*sql.Rows, error)
}
