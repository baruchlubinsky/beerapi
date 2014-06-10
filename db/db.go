package db

type Model interface {
	Id() (Id, error)
	Save() (bool, error)
	Marshal() ([]byte, error)
}

type Table interface {
	Create(Attributes) (*Model, error)
	Find(Id) (*Model, error)
	Search(...interface{}) ([]Model, error)
}

type Attributes map[string]interface{}

type Id string