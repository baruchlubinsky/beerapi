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
	Db.CreateTable("beers")
	Db.CreateTable("comments")
}

func beer(response http.ResponseWriter, request *http.Request) {
	header := response.Header()
	header.Set("Content-Type", "application/json")
	// CORS
	header.Add("Access-Control-Allow-Origin","*")
	header.Add("Access-Control-Allow-Methods","POST, PUT, DELETE, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, content-type, Accept, X-AUTH-TOKEN, X-API-VERSION")
	// Check for table 
	table, err := tableFor(request)
	if err != nil {
		response.Write([]byte(err.Error()))
		response.WriteHeader(404)
		return
	}
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

func tableFor(request *http.Request) (adapters.Table, error) {
	args := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	name := args[0]
	table := Db.Table(name)
	if table == nil {
		return nil, db.DBError("Table not found: " + name)
	}
	return table, nil
}

const PORT = ":9000"

func main () {
	log.Println("Connected to database. Listening on " + PORT)
	http.HandleFunc("/", beer)
	http.ListenAndServe(PORT, nil)
}