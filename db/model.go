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

func (model *Model) SetId() Id {
	hash := strconv.Itoa(int(time.Now().Unix())) + fmt.Sprint(model)
	raw := sha256.Sum256([]byte(hash))
	model.Id = Id(fmt.Sprintf("%v", hex.EncodeToString(raw[:16])))
	model.data["id"] = model.Id
	return model.Id
}

func (model *Model) Attributes() Attributes {
	return model.data
}

func (model *Model) SetAttributes(attributes Attributes) {
	for key, value := range(attributes) {
		model.data[key] = value
	}
}

func (model *Model) Marshal(name string) ([]byte, error) {
	data := map[string]interface{}{name: model.data}
	return json.Marshal(data)
}

func (model *Model) Save() (error) {
	return model.table.Save(model)
}