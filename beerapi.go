package main

import (
	"github.com/crowdmob/goamz/dynamodb"
	"github.com/crowdmob/goamz/aws"
	_"net/http"
	"fmt"
	"time"
	"io/ioutil"
	"strings"
)

var dbserver dynamodb.Server

var auth aws.Auth

var region = aws.Region {
	Name: "eu-west-1",
	DynamoDBEndpoint: "http://dynamodb.eu-west-1.amazonaws.com",
}

func init () {
	file, _ := ioutil.ReadFile("credentials.csv")
	lines := strings.Split(string(file), "\n")
	data := stirng.Split(lines[1], ",")
	auth, err := aws.GetAuth(data[1], data[2], "", time.Now())
	if err != nil {
		panic("Unable to connect to aws.")
	}
	dbserver = dynamodb.Server {auth, region}
}

func main () {
	fmt.Println("Connected to database.")

	beerTable := dynamodb.Table {
		&dbserver,
		"Beer",
		dynamodb.PrimaryKey{dynamodb.NewStringAttribute("Name", ""),nil},
	}

	res, err := beerTable.PutItem("Castle", "", []dynamodb.Attribute{*dynamodb.NewStringAttribute("Type","Lager")})
	if res {
		fmt.Println("Created item.")
	}
	if err != nil {
		panic(err)
	}
}