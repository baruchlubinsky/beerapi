package main

import (
	"github.com/crowdmob/goamz/dynamodb"
	"github.com/crowdmob/goamz/aws"
	"net/http"
	"log"
	"time"
	"io/ioutil"
	"strings"
	"encoding/json"
)

var dbserver dynamodb.Server

var auth aws.Auth

var beerTable dynamodb.Table

var region = aws.Region {
	Name: "eu-west-1",
	DynamoDBEndpoint: "http://dynamodb.eu-west-1.amazonaws.com",
}

type Beer struct {
	Name string
	Type string
}

func init () {
	file, _ := ioutil.ReadFile("credentials.csv")
	lines := strings.Split(string(file), "\n")
	data := strings.Split(lines[1], "\t")
	auth, err := aws.GetAuth(data[1], data[2], "", time.Now())
	if err != nil {
		panic("Unable to connect to aws.")
	}
	dbserver = dynamodb.Server {auth, region}

	beerTable = dynamodb.Table {
		&dbserver,
		"Beer",
		dynamodb.PrimaryKey{dynamodb.NewStringAttribute("Name", ""),nil},
	}
}

func create(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		body, err := ioutil.ReadAll(request.Body)
		check(err)
		var data Beer
		err = json.Unmarshal(body, &data) 
		check(err)
		log.Printf("POST /beer %#v", data)
		_, err = beerTable.PutItem(data.Name, "", []dynamodb.Attribute{*dynamodb.NewStringAttribute("Type",data.Type)})
		check(err)
		response.WriteHeader(http.StatusOK)	
	} else {
		response.WriteHeader(400)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main () {
	log.Println("Connected to database.")

	http.HandleFunc("/beer", create)
	http.ListenAndServe(":8080", nil)
}