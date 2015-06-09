package api

import (
	"github.com/baruchlubinsky/beerapi/adapters"
	"net/http"
	"strings"
)

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
		var data adapters.ModelSet
		if ids, q := query["ids[]"]; q {
			data = make(adapters.ModelSet, 0)
			for _, id := range ids {
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
