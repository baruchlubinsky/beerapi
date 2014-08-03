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

func NewTable(name string) *Table {
	res := Table {
		data: make([]*Model, 0),
		index: make(map[Id]int),
		plural: name,
		singular: name[0:len(name) - 1],
	}
	return &res
}

func (table *Table) Save(model *Model) (err error) {
	if model.Id == "" {
		model.SetId()
	}
	index, found := table.index[model.Id]
	if found {
		table.data[index] = model
	} else {
		table.data = append(table.data, model)
		// if table.index == nil {
		// 	table.index = make(map[Id]int)
		// }
		table.index[model.Id] = len(table.data) - 1
	} 
    return
}

func (table *Table) Find(id string) (model adapters.Model, err error) {
	index, found := table.index[Id(id)]
	if found {
		return table.data[index], nil
	} else {
		return nil, DBError("Object with that ID does not exist.")
	}
}

func (table *Table) Search(query interface{}) (result adapters.ModelSet) {
	result = make(adapters.ModelSet, 0)
	for _, model := range table.All() {
		match := query == nil
		for key, value := range query.(map[string]interface{}) {
			if !reflect.DeepEqual(model.(*Model).data[key], value) {
				break
			}
			match = true
		}
		if match {
			result.Add(model)
		}
	}
	return
}

func (table *Table) All() (result adapters.ModelSet) {
	result = make(adapters.ModelSet, 0)
	for _, index := range table.index {
		result.Add(table.data[index])
	}
	return
}

func (table *Table) NewRecord() (adapters.Model) {
	return &Model{
		table: table,
		data: make(map[string]interface{}),
	}
}

func (table *Table) Delete(id string) (error) {
	if index, found := table.index[Id(id)]; found {
		table.data[index] = nil
		delete(table.index, Id(id))
		return nil
	} else {
		return DBError("Object with that ID does not exist.")
	}
}

func (table *Table) RecordName() string {
	return table.singular
}

func (table *Table) RecordSetName() string {
	return table.plural
}