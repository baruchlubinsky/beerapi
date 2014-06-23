package db

import (
	"time"
	"crypto/sha256"
	"strconv"
	"fmt"
	"encoding/json"
	"encoding/hex"
)

type Model struct {
	Id Id
	data Attributes
	table *Table
}

type ModelSet []*Model

func (model *Model) SetId() string {
	hash := strconv.Itoa(int(time.Now().Unix())) + fmt.Sprint(model)
	raw := sha256.Sum256([]byte(hash))
	model.Id = Id(fmt.Sprintf("%v", hex.EncodeToString(raw[:16])))
	model.data["id"] = model.Id
	return string(model.Id)
}

func (model *Model) Attributes() interface{} {
	return model.data
}

func (model *Model) SetAttributes(attributes interface{}) {
	for key, value := range(attributes.(Attributes)) {
		model.data[key] = value
	}
}

func (model *Model) Marshal(name string) ([]byte, error) {
	data := map[string]interface{}{name: model.data}
	return json.Marshal(data)
}

func (set ModelSet) Marshal(name string) ([]byte, error) {
	rows := make([]interface{}, len(set))
	for i, model := range set {
		rows[i] = model.data
	}
	data := map[string]interface{}{name: rows}
	return json.Marshal(data)
}

func (model *Model) Save() (error) {
	return model.table.Save(model)
}