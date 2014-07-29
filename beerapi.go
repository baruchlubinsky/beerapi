package main

import (
	"net/http"
	
	"io/ioutil"
	"encoding/json"
	
	"github.com/baruchlubinsky/beerapi/db"
	"strings"
	//"github.com/gedex/inflector"
	"log"
	
	
)

// type Beer struct {
// 	id db.Id `json:"id"`
// 	Name string `json:"name"`
// 	Type string `json:"type"`
// 	Comments []db.Id `json:"comments"`
// }

// type Comment struct {
// 	Id db.Id `json:"id"`
// 	Text string `json:"text"`
// 	BeerId db.Id `json:"beer"`
// }

// type CommentTable struct {
// 	data []*Comment
// 	index map[db.Id]int
// }

type Marshalable interface{
	Marshal(string) ([]byte, error)
}

func Unmarshal(data []byte, name string) (db.Attributes, error) {
	var object map[string]interface{}
	err := json.Unmarshal(data, &object)
	return db.Attributes(object[name].(map[string]interface{})), err
}

// Database interface 

type Id string

type Attributes db.Attributes

type DatabaseAdapter interface {
	CreateTable(string)
	Table(string) *TableAdapter
}

type TableAdapter interface {
	Find(id string) (ModelAdapter, error)
	Search(query Attributes) (result ModelSetAdapter)
	NewRecord() (ModelAdapter)
}

type ModelAdapter interface {
	SetId() string
	Attributes() interface{}
	SetAttributes(interface{})
	Marshal(string) ([]byte, error)
	Save() (error)
}

type ModelSetAdapter interface {
	Marshal(string) ([]byte, error)
}

var Db db.Database

func init() {
	//Db = make(map[string]*db.Table)
	// Db["beers"] = db.NewTable()
	Db.CreateTable("beers")
	Db.CreateTable("comments")
}

func get(response http.ResponseWriter, request *http.Request) {
	args := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	table := Db.Table(args[0])
	var data interface{}
	var name string
	if table == nil {
		log.Println("Table not found.")
		response.WriteHeader(404)
		return
	}
	if len(args) == 1 {
		query := request.URL.Query()
		ids, q := query["ids[]"]
		if q {
			data = make([]ModelAdapter, 0)
			for _, id := range(ids) {
				record, err := table.Find(id)
				if err != nil {
					response.WriteHeader(404)
					return
				}
				data = append(data.([]ModelAdapter), record)
			}
		} else {
			data = table.Search(nil)
		}
		name = args[0]
	} else if len(args) == 2 {
		id := args[1]
		record, err := table.Find(id)
		if err != nil {
			response.WriteHeader(404)
			return
		}
		data = record
		name = args[0][:len(args[0])]
	} else {
		response.WriteHeader(404)
		return
	}
	resp, _ := data.(Marshalable).Marshal(name)
	response.Write(resp)
}

func create(response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	check(err)
	log.Println(string(data))
	object, err := Unmarshal(data, "beer")
	check(err)
	table := Db.Table("beers")
	if table == nil {
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

func put(response http.ResponseWriter, request *http.Request) {
	log.Println("called put")
	args := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	name := args[0]
	table := Db.Table(name)
	if table == nil {
		log.Println("Table not found.")
		response.WriteHeader(404)
		return
	}
	data, err := ioutil.ReadAll(request.Body)
	check(err)
	log.Println(string(data))
	object, err := Unmarshal(data, name[0:len(name) - 1])
	check(err)
	if len(args) == 2 {
		id := args[1]
		record, err := table.Find(id)
		if err != nil {
			response.WriteHeader(404)
			return
		}
		record.SetAttributes(object)
		record.Save()
		resp, _ := record.Marshal(name[0:len(name) - 1])
		response.Write(resp)
	} else {
		log.Println("No id")
		response.WriteHeader(500)
		return
	}
}

func beer(response http.ResponseWriter, request *http.Request) {
	header := response.Header()
	header.Set("Content-Type", "application/json")
	// CORS
	header.Add("Access-Control-Allow-Origin","*")
	header.Add("Access-Control-Allow-Methods","POST, PUT, DELETE, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, content-type, Accept, X-AUTH-TOKEN, X-API-VERSION")
	args := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	log.Printf("%v: %#v\n", request.Method, args)
	switch request.Method {
	case "POST": 
		create(response, request)
	case "GET":
		get(response, request) 
	case "PUT":
		put(response, request)
	case "OPTIONS":
		response.WriteHeader(200)
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

const PORT = ":9000"

func main () {
	log.Println("Connected to database. Listening on " + PORT)
	http.HandleFunc("/", beer)
	http.ListenAndServe(PORT, nil)
}