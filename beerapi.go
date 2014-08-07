// This program is intended only as an example of the usage of the 
// beerapi/api package.
package main

import (
	"net/http"	
	"github.com/baruchlubinsky/beerapi/adapters"
	"github.com/baruchlubinsky/beerapi/api"
	"github.com/baruchlubinsky/beerapi/db"
	"strings"
	"log"	
)

var Db adapters.Database

func init() {
	Db = &db.Database{}
}

func beer(response http.ResponseWriter, request *http.Request) {
	header := response.Header()
	header.Set("Content-Type", "application/json")
	// CORS
	header.Add("Access-Control-Allow-Origin","*")
	header.Add("Access-Control-Allow-Methods","POST, PUT, DELETE, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, content-type, Accept, X-AUTH-TOKEN, X-API-VERSION")
	// Check for table 
	table := tableFor(request)
	switch request.Method {
	case "POST": 
		api.Post(table, response, request)
	case "GET":
		api.Get(table, response, request) 
	case "PUT":
		api.Put(table, response, request)
	case "DELETE":
		api.Delete(table, response, request)
	case "OPTIONS":
		response.WriteHeader(200)
	default:
		response.WriteHeader(400)
	}
}

func tableFor(request *http.Request) (adapters.Table) {
	args := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	name := args[0]
	table, err := Db.Table(name)
	if err != nil {
		return Db.CreateTable(name)
	}
	return table
}

const PORT = ":9000"

func main () {
	log.Println("Connected to database. Listening on " + PORT)
	http.HandleFunc("/", beer)
	http.ListenAndServe(PORT, nil)
}