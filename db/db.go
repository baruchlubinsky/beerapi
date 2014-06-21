package db

type Database map[string]*Table

type Attributes map[string]interface{}

type Id string

func (database Database) CreateTable(name string) {
	//map[string]*Table(database)[name] = NewTable()
	database[name] = NewTable()
}

type DBError string

func (a DBError) Error() string {
	return string(a)
}