package db

import (
	"time"
	"crypto/sha256"
	"strconv"
	"fmt"
	"encoding/hex"
)

type Model struct {
	Id Id
	data map[string]interface{}
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
	for key, value := range(attributes.(map[string]interface{})) {
		model.data[key] = value
	}
}

func (model *Model) Save() (error) {
	return model.table.Save(model)
}

func (model *Model) Delete() (error) {
	return model.table.Delete(string(model.Id))
}