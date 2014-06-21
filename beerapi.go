package main

import (
	"net/http"
	
	"io/ioutil"
	"encoding/json"
	
	"github.com/baruchlubinsky/beerapi/db"
	
	
	"log"
	
	
)


type Beer struct {
	id db.Id `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Comments []db.Id `json:"comments"`
}

type Comment struct {
	Id db.Id `json:"id"`
	Text string `json:"text"`
	BeerId db.Id `json:"beer"`
}

type CommentTable struct {
	data []*Comment
	index map[db.Id]int
}

// Model Interface

func Unmarshal(data []byte) (db.Attributes, error) {
	var object db.Attributes
	err := json.Unmarshal(data, &object)
	return object, err
}

// Table Interface

var Db map[string]*db.Table

func init() {
	Db = make(map[string]*db.Table)
	//Db["beers"] = db.NewTable()
	db.Database(Db).CreateTable("beers")
}

func get(response http.ResponseWriter, request *http.Request) {
	id := request.URL.Path[len("/beers/"):]
	table, found := Db["beers"]
	if !found {
		response.WriteHeader(404)
		return
	}
	record, err := table.Find(db.Id(id))
	if err != nil {
		response.WriteHeader(401)
	} else {
		resp, _ := record.Marshal("beer")
		response.Write(resp)
	}
}

func create(response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	check(err)
	object, err := Unmarshal(data)
	check(err)
	table, found := Db["beers"]
	if !found {
		response.WriteHeader(404)
		return
	}
	record := table.NewRecord()
	record.SetAttributes(object)
	table.Save(record)
	if err != nil {
		response.WriteHeader(401)
	} else {
		resp, err := record.Marshal("beer")
		check(err)
		response.Write(resp)
	}
}

func beer(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST": 
		create(response, request)
	case "GET":
		get(response, request) 
	default:
		response.WriteHeader(400)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type ApiError string

func (a ApiError) Error() string {
	return string(a)
}

const PORT = ":8080"

func main () {
	log.Println("Connected to database. Listening on " + PORT)

	http.HandleFunc("/beers/", beer)
	http.ListenAndServe(PORT, nil)
}