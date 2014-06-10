package main

import (
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	"crypto/sha256"
	"github.com/baruchlubinsky/beerapi/db"
	"fmt"
	"strconv"
	"log"
)


type Beer struct {
	Id db.Id `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type BeerTable struct {
	data []*Beer
	index map[db.Id]int
}

// Model Interface

func Unmarshal(data []byte) (db.Attributes, error) {
	var object db.Attributes
	err := json.Unmarshal(data, &object)
	return object, err
}

func (beer *Beer) SetId() db.Id {
	hash := strconv.Itoa(int(time.Now().Unix())) + beer.Name + beer.Type
	raw := sha256.Sum256([]byte(hash))
	beer.Id = db.Id(fmt.Sprintf("%v", raw))
	return beer.Id
}

func (beer *Beer) Marshal() ([]byte, error) {
	data := struct{
		Beer Beer `json:"beer"`
	}{
		Beer: *beer,
	}
	return json.Marshal(data)
}

// Table Interface

func (table *BeerTable) Init() {
	table.data = make([]*Beer, 0, 10)
	table.index = make(map[db.Id]int, 10)
}

func (table *BeerTable) Create(data db.Attributes) (beer *Beer, err error) {
	fmt.Println(data)
	beer = &Beer{
		Name: data["name"].(string),
    	Type: data["type"].(string),
    }
	table.data = append(table.data, beer)
    table.index[beer.SetId()] = len(table.data) - 1
    return
}

func (table *BeerTable) Find(id db.Id) (beer *Beer, err error) {
	index, found := table.index[id]
	if found {
		return table.data[index], nil
	} else {
		return nil, ApiError("Object with that ID does not exist.")
	}
}


var Db BeerTable

func init() {
	Db.Init()
}

func get(response http.ResponseWriter, request *http.Request) {
	id := request.URL.Path[len("/beers/"):]
	beer, err := Db.Find(db.Id(id))
	if err != nil {
		response.WriteHeader(401)
	} else {
		resp, _ := beer.Marshal()
		response.Write(resp)
	}
}

func create(response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	check(err)
	object, err := Unmarshal(data)
	check(err)
	beer, err := Db.Create(object)
	if err != nil {
		response.WriteHeader(401)
	} else {
		resp, _ := beer.Marshal()
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