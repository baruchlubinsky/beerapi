package api

import (
	"github.com/baruchlubinsky/beerapi/adapters"
	"io/ioutil"
	"net/http"
)

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
