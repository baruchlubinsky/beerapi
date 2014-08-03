package adapters

type Database interface {
	CreateTable(string)
	Table(string) Table
}

type Table interface {
	Find(string) (Model, error)
	Search(query interface{}) (result ModelSet)
	NewRecord() (Model)
	Delete(string) (error)
	RecordName() (string)
	RecordSetName() (string)
}

type Model interface {
	SetId() string
	Attributes() interface{}
	SetAttributes(interface{})
	Save() (error)
	Delete() (error)
}

type ModelSet []Model 

func (set ModelSet) Add(model Model) {
	set = append(set, model)
}