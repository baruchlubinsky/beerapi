package db

import (
	"time"
	"crypto/sha256"
	"strconv"
	"fmt"
	"encoding/hex"
)

type Model struct {
	data map[string]interface{}
	table *Table
}

type ModelSet []*Model

// Returns and lazily creates this model's ID.
func (model *Model) GetId() string {
	id, set := model.data["id"]
	if !set {
		id = model.setId()
	}
	return id.(string)
}

func (model *Model) setId() string {
	hash := strconv.Itoa(int(time.Now().Unix())) + fmt.Sprint(model)
	raw := sha256.Sum256([]byte(hash))
	id := fmt.Sprintf("%v", hex.EncodeToString(raw[:16]))
	model.data["id"] = id
	return id
}

// Return this model's data. 
func (model *Model) Attributes() interface{} {
	return model.data
}

// Set the model's data. In this implementation, expects attribures to 
// be map[string]interface{}.
func (model *Model) SetAttributes(attributes interface{}) {
	for key, value := range(attributes.(map[string]interface{})) {
		model.data[key] = value
	}
}

// Store this model in the database.
func (model *Model) Save() (error) {
	// Since table.data is []*Model, this is redundant.
	return model.table.Save(model)
}

// Delete this model from the database.
func (model *Model) Delete() (error) {
	return model.table.Delete(model.GetId())
}