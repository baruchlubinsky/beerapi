package db 

import (
	"reflect"
	"github.com/baruchlubinsky/beerapi/adapters"
)

type Table struct {
	data []*Model
	index map[Id]int
	singular, plural string
}

// Create a new table with specified name. Name should be a plural, singular
// form is assumed to be name minus its last letter.
func NewTable(name string) *Table {
	res := Table {
		data: make([]*Model, 0),
		index: make(map[Id]int),
		plural: name,
		singular: name[0:len(name) - 1],
	}
	return &res
}

// Create or update the specified record.
func (table *Table) Save(model *Model) (err error) {
	id := Id(model.GetId())
	index, found := table.index[id]
	if found {
		table.data[index] = model
	} else {
		table.data = append(table.data, model)
		table.index[id] = len(table.data) - 1
	} 
    return
}

// Returns the object with specified ID. If the ID is not found error is
// DBError and model is nil.
func (table *Table) Find(id string) (model adapters.Model, err error) {
	index, found := table.index[Id(id)]
	if found {
		return table.data[index], nil
	} else {
		return nil, DBError("Object with that ID does not exist.")
	}
}

// Return the records which match query, or all records in a table when query
// is nil. In this impelementation query is expected to be map[string]interface.
// Records are returned if all the fields contained in query are equal by
// reflect.DeepEqual to their counterparts in the record.
func (table *Table) Search(query interface{}) (result adapters.ModelSet) {
	q, provided := query.(map[string]interface{})
	if provided {
		result = make(adapters.ModelSet, 0)
		for _, model := range table.All() {
			match := false
			for key, value := range q {
				if !reflect.DeepEqual(model.(*Model).data[key], value) {
					break
				}
				match = true
			}
			if match {
				result.Add(model)
			}
		}
	} else {
		result = table.All()
	}
	return
}

// Returns all the records in this table.
func (table *Table) All() (result adapters.ModelSet) {
	result = make(adapters.ModelSet, 0)
	for _, index := range table.index {
		result.Add(table.data[index])
	}
	return
}

// Create a new record for this table. The new record is not added to the
// table until Save is called.
func (table *Table) NewRecord() (adapters.Model) {
	return &Model{
		table: table,
		data: make(map[string]interface{}),
	}
}

// Remove the record with specified ID from the table. Returns a DBError id
// that ID does not exist.
func (table *Table) Delete(id string) (error) {
	if index, found := table.index[Id(id)]; found {
		table.data[index] = nil
		delete(table.index, Id(id))
		return nil
	} else {
		return DBError("Object with that ID does not exist.")
	}
}

// The name of an individual record.
func (table *Table) RecordName() string {
	return table.singular
}

// The name of a collection of records.
func (table *Table) RecordSetName() string {
	return table.plural
}