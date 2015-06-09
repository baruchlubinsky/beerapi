package api

import (
	"github.com/baruchlubinsky/beerapi/adapters"
	"net/http"
)

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
