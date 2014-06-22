package db

type Database struct{
	tables map[string]*Table
}

type Attributes map[string]interface{}

type Id string

func (database *Database) CreateTable(name string) {
	if database.tables == nil {
		database.tables = make(map[string]*Table)
	}
	database.tables[name] = NewTable()
}

func (database *Database) Table(name string) *Table {
	return database.tables[name]
}

type DBError string

func (a DBError) Error() string {
	return string(a)
}