This webserver is a backend for [Ember Data](http://emberjs.com/api/data/). It is intended to be used as documentation of the structure of JSON that Ember Data expects and as an alternative to using fixtures during development. 

Usage
==

The command `beerapi` is a standard [Go webservice](http://golang.org/doc/articles/wiki/). 

The tables that the server will handle must be defined in `init`. (This could easily be changed.)

The `db` package in this repository is simply an example implimentation of the interfaces in `adapters`. It has no persistence. To use a different database, replace the line `Db = &db.Database{}`.

Find the Ember application demonstrating Ember Data which is designed together with this backend at baruchlubinsky/beerdemo .

App Engine
==

Use this package on [App Engine](https://developers.google.com/appengine/docs/go/gettingstarted/introduction).

Simply copy beerapi.go, move `http.HandleFunc("/", beer)` into `init()` and delete `main()`.

Create a `dispatch.yaml` with:

	application: restapi

	dispatch:
	    - url: "*/*"
	      module: "default"

And `app.yaml`:

	application: restapi

		version: 1
		runtime: go
		api_version: go1

		module: default
		handlers:
		- url: /.*
		  script: _go_app

