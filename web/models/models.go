package models

import (
	"database/sql"
)
import _ "github.com/go-sql-driver/mysql"

type Model struct {
	dbSource string
}

func New(dbSource string) *Model {
	return &Model{
		dbSource: dbSource,
	}
}

func (m *Model) Open() (*sql.DB, error) {
	return sql.Open("mysql", m.dbSource)
}
