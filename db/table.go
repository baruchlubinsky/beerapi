package db 

import (
	"reflect"
)

type Table struct {
	data []*Model
	index map[Id]int
}

func NewTable() *Table {
	res := Table {
		data: make([]*Model, 0),
		index: make(map[Id]int),
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

func (table *Table) Find(id string) (model *Model, err error) {
	index, found := table.index[Id(id)]
	if found {
		return table.data[index], nil
	} else {
		return nil, DBError("Object with that ID does not exist.")
	}
}

func (table *Table) Search(query Attributes) (result ModelSet) {
	for _, model := range table.data {
		match := query == nil
		for key, value := range query {
			if !reflect.DeepEqual(model.Attributes().(Attributes)[key], value) {
				break
			}
			match = true
		}
		if match {
			result = append(result, model)
		}
	}
	return
}

func (table *Table) NewRecord() (*Model) {
	return &Model{
		table: table,
		data: make(map[string]interface{}),
	}
}