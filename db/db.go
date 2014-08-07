// Package db implements a database as required in beerapi/adapters. 
// It does not persist its data, all records are stored in memory. It 
// is not intended to be used in and sort of production environment, 
// but rather as an example of how to implement the adapters interface.
package db

import (
	"github.com/baruchlubinsky/beerapi/adapters"
)

type Database struct{
	tables map[string]*Table
}

type Id string

// Create a table in the database with specified name.
func (database *Database) CreateTable(name string) (adapters.Table) {
	if database.tables == nil {
		database.tables = make(map[string]*Table)
	}
	database.tables[name] = NewTable(name)
	return database.tables[name]
}

// Get the table with specified name, returns nil if that table does not exist.
func (database *Database) Table(name string) (adapters.Table, error) {
	table, found := database.tables[name]
	if found {
		return table, nil
	}
	return nil, DBError("Table not found.")
}

type DBError string

func (a DBError) Error() string {
	return string(a)
}