// Package api implements the REST verbs. It also contains functions for
// serializing and deserialing models as JSON. This pacakge is independant
// of the implementation of the database adapter.
package api

import (
	"encoding/json"
	"github.com/baruchlubinsky/beerapi/adapters"
	"net/http"
	"strings"
)

type HandlerFunc func(adapters.Table, http.ResponseWriter, *http.Request)

// Returns the JSON represenation of this model with root element name.
func Marshal(model adapters.Model, name string) ([]byte, error) {
	data := map[string]interface{}{name: model.Attributes()}
	return json.Marshal(data)
}

// Returns the JSON representation of this slice of models, with root
// element name which should be a plural.
func MarshalSet(set adapters.ModelSet, name string) ([]byte, error) {
	rows := make([]interface{}, len(set))
	for i, model := range set {
		rows[i] = model.Attributes()
	}
	data := map[string]interface{}{name: rows}
	return json.Marshal(data)
}

// Deserialize JSON data from a request, name should be the root element.
func Unmarshal(data []byte, name string) (map[string]interface{}, error) {
	var object map[string]interface{}
	err := json.Unmarshal(data, &object)
	return object[name].(map[string]interface{}), err
}

func getArgs(request *http.Request) []string {
	return strings.Split(strings.Trim(request.URL.Path, "/"), "/")
}

func write(response http.ResponseWriter, data interface{}, table adapters.Table) {
	var resp []byte
	var err error
	switch data.(type) {
	case adapters.Model:
		resp, err = Marshal(data.(adapters.Model), table.RecordName())
	case adapters.ModelSet:
		resp, err = MarshalSet(data.(adapters.ModelSet), table.RecordSetName())
	default:
		panic("Attempt to write unknown type as response.")
	}
	check(err)
	response.Write(resp)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
