// Package api implements the REST verbs. It also conatains functions for 
// serializing and deserialing models. This pacakge is independant of the
// implementation of the database adapter.
package api

import (
	"net/http"	
	"encoding/json"
	"io/ioutil"
	"github.com/baruchlubinsky/beerapi/adapters"
	"strings"
)

type HandlerFunc func(adapters.Table, http.ResponseWriter, *http.Request)

// Returns the JSON represenation of this model with root element `name`.
func Marshal(model adapters.Model, name string) ([]byte, error) {
	data := map[string]interface{}{name: model.Attributes()}
	return json.Marshal(data)
}

// Returns the JSON representation of this slice of models, with root 
// element `name` which should be a plural.
func MarshalSet(set adapters.ModelSet, name string) ([]byte, error) {
	rows := make([]interface{}, len(set))
	for i, model := range set {
		rows[i] = model.Attributes()
	}
	data := map[string]interface{}{name: rows}
	return json.Marshal(data)
}

// Deserialize JSON data from a request, `name` should be the root element.
func Unmarshal(data []byte, name string) (map[string]interface{}, error) {
	var object map[string]interface{}
	err := json.Unmarshal(data, &object)
	return object[name].(map[string]interface{}), err
}

func Get(dbTable adapters.Table, response http.ResponseWriter, request *http.Request) {
	args := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	if len(args) == 1 {
		query := request.URL.Query()
		ids, q := query["ids[]"]
		var data adapters.ModelSet
		if q {
			data = make(adapters.ModelSet, 0)
			for _, id := range(ids) {
				record, err := dbTable.Find(id)
				if err != nil {
					response.WriteHeader(404)
					return
				}
				data.Add(record)
			}
		} else {
			data = dbTable.Search(nil)
		}
		resp, err := MarshalSet(data, dbTable.RecordSetName())
		check(err)
		response.Write(resp)
	} else if len(args) == 2 {
		id := args[1]
		record, err := dbTable.Find(id)
		if err != nil {
			response.WriteHeader(404)
		} else {
			resp, err := Marshal(record, dbTable.RecordName())
			check(err)
			response.Write(resp)
		}
	} else {
		response.WriteHeader(404)
	}
}

func Create(dbTable adapters.Table, response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	check(err)
	object, err := Unmarshal(data, dbTable.RecordName())
	check(err)
	record := dbTable.NewRecord()
	record.SetAttributes(object)
	err = record.Save()
	if err != nil {
		response.WriteHeader(401)
	} else {
		resp, err := Marshal(record, dbTable.RecordName())
		check(err)
		response.Write(resp)
	}
}

func Put(dbTable adapters.Table, response http.ResponseWriter, request *http.Request) {
	args := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	data, err := ioutil.ReadAll(request.Body)
	check(err)
	object, err := Unmarshal(data, dbTable.RecordName())
	check(err)
	if len(args) == 2 {
		id := args[1]
		record, err := dbTable.Find(id)
		if err != nil {
			response.WriteHeader(404)
			return
		}
		record.SetAttributes(object)
		err = record.Save()
		if err != nil {
			response.WriteHeader(401)
		} else {
			resp, err := Marshal(record, dbTable.RecordName())
			check(err)
			response.Write(resp)
		}
	} else {
		response.Write([]byte("No id"))
		response.WriteHeader(500)
	}
}

func Delete(dbTable adapters.Table, response http.ResponseWriter, request *http.Request) {
	args := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	if len(args) == 2 {
		id := args[1]
		record, err := dbTable.Find(id)
		if err != nil {
			response.WriteHeader(404)
		} else {
			err = record.Delete()
			check(err)
			response.Write([]byte("{}"))
		}
	} else {
		response.Write([]byte("No id"))
		response.WriteHeader(500)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}