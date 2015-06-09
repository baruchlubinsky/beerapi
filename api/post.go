package api

import (
	"github.com/baruchlubinsky/beerapi/adapters"
	"io/ioutil"
	"net/http"
)

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
