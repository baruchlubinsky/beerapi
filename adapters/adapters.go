// Package adapters contains the interfaces required to define a general
// database adapter. See beerapi/db for documentation.
package adapters

type Database interface {
	CreateTable(string) Table
	Table(string) (Table, error)
}

type Table interface {
	Find(string) (Model, error)
	// If query == nil, return entire contents of table.
	Search(query interface{}) (result ModelSet)
	NewRecord() Model
	Delete(string) error
	RecordName() string
	RecordSetName() string
}

type Model interface {
	GetId() string
	Attributes() map[string]interface{}
	SetAttributes(map[string]interface{})
	Save() error
	Delete() error
}

type ModelSet []Model

func (set *ModelSet) Add(model Model) {
	s := *set
	s = append(s, model)
	*set = s
}
