// Package api implements the REST verbs. It also contains functions for 
// serializing and deserialing models as JSON. This pacakge is independant
// of the implementation of the database adapter.
package api

import (
	"net/http"	
	"encoding/json"
	"io/ioutil"
	"github.com/baruchlubinsky/beerapi/adapters"
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

// Handle a GET request. This function supports three options:
//    - GET /beers/ - returns all the records in the beers table.
//    - GET /beers/:id - returns the beer with specified id
//    - GET /beers?ids[]=1&ids[]=2 - returns all the specified beer records.
//         see http://emberjs.com/api/data/classes/DS.RESTAdapter.html#method_findMany
//
// If an id is specified but not found, returns a 404 error.
//
// If none of the above patterns is follwed, return a 400 error.
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
		write(response, data, dbTable)
	} else if len(args) == 2 {
		id := args[1]
		record, err := dbTable.Find(id)
		if err != nil {
			response.WriteHeader(404)
		} else {
			write(response, record, dbTable)
		}
	} else {
		response.WriteHeader(400)
	}
}

// Handle a POST request.
//
// If the body of the request cannot be deserialized, panic.
//
// If saving the data returns an error, return a 400 error. 
//
// Returns the saved record.
func Post(dbTable adapters.Table, response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	check(err)
	object, err := Unmarshal(data, dbTable.RecordName())
	check(err)
	record := dbTable.NewRecord()
	record.SetAttributes(object)
	err = record.Save()
	if err != nil {
		response.Write([]byte(err.Error()))
		response.WriteHeader(400)
	} else {
		write(response, record, dbTable)
	}
}

// Handle a PUT request. 
//    - PUT /beers/:id
// 
// If no id is specified returns a 400 error.
//
// If the specified id is not found, returns a 404 error.
// 
// If the body of the request cannot be deserialized, panic.
//
// If saving the data returns an error, return a 400 error. 
//
// Returns the saved record.
func Put(dbTable adapters.Table, response http.ResponseWriter, request *http.Request) {
	args := getArgs(request)
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
			response.Write([]byte(err.Error()))
			response.WriteHeader(400)
		} else {
			write(response, record, dbTable)
		}
	} else {
		response.Write([]byte("No id"))
		response.WriteHeader(400)
	}
}

// Handle a DELETE request. 
//    - DELETE /beers/:id
// 
// If no id is specified returns a 400 error.
//
// If the specified id is not found, returns a 404 error.
// 
// If saving the data returns an error, return a 400 error. 
//
// Returns "{}" with code 200.
func Delete(dbTable adapters.Table, response http.ResponseWriter, request *http.Request) {
	args := getArgs(request)
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
		response.WriteHeader(400)
	}
}

func getArgs(request *http.Request) ([]string) {
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