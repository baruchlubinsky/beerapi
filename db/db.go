package db

import (
	"github.com/baruchlubinsky/beerapi/adapters"
)

type Database struct{
	tables map[string]*Table
}

type Id string

func (database *Database) CreateTable(name string) {
	if database.tables == nil {
		database.tables = make(map[string]*Table)
	}
	database.tables[name] = NewTable(name)
}

func (database *Database) Table(name string) adapters.Table {
	return database.tables[name]
}

type DBError string

func (a DBError) Error() string {
	return string(a)
}