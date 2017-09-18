// Test the library
package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"time"

	"github.com/florenthobein/godoc2api"
)

func main() {

	logLevel := flag.Uint("log", 3, "Define a log level") // default: nothing
	flag.Parse()
	if logLevel != nil && *logLevel < 3 {
		godoc2api.LogLevel = *logLevel
	}

	output_dir := "output"

	// Configuration
	godoc2api.DefineSecurity("auth", nil)    // todo
	godoc2api.DefineTrait("deprecated", nil) // todo
	godoc2api.DefineType("MyStruct", MyStruct{})
	godoc2api.DefineType("MyStruct2", MyStruct2{})

	// Create the doc
	doc := godoc2api.Documentation{
		Title:       "Test API",
		Description: "API used for tests",
		Version:     "v1",
		URL:         "http://mywebsite/{version}",
	}

	// Add the a full route definition, the handler is not commented
	doc.AddRoute(RouteDefinition{
		Method:      "POST",
		Resource:    "/myroute/{id}",
		Description: "A route that use a handler without comments",
		Handler:     MyHanderWithoutComment,
		RouteParams: [][]string{[]string{"{string}", "id", "The id of my route"}},
		QueryParams: [][]string{[]string{"{bool}", "[working=true]", "If set to `true`, everything works just fine"}},
		Body:        "{MyStruct2}",
		Examples: [][]string{
			[]string{"A complicated test", "{ \"value_6\": { \"test\": true } }", "200: {" +
				"\"value_1\": \"\", \"value_2\": 0, \"value_3\": false" +
				"}"},
			[]string{"A simpler test", "/myroute/1?working=false", "200: nil"},
		},
		Response: "{MyStruct}",
	})

	// Add a route with just the handler
	doc.AddRoute(MyHanderWithAllTheComments)

	// Add a route definition completed with comments
	doc.AddRoute(RouteDefinition{
		Resource: "GET /myroute/{id}",
		Handler:  MyHanderWithFewComments,
		Auth:     true,
	})

	// Render the RAML document
	doc.Render(output_dir)
}

type RouteDefinition struct {
	Method      string                                   `raml:"method"`
	Resource    string                                   `raml:"resource"`
	Description string                                   `raml:"description"`
	Handler     func(http.ResponseWriter, *http.Request) `raml:"handler"`
	RouteParams [][]string                               `raml:"routes"`
	QueryParams [][]string                               `raml:"queries"`
	Body        string                                   `raml:"body"`
	Response    string                                   `raml:"response"`
	Examples    [][]string                               `raml:"examples"`
	Auth        bool                                     `raml:"auth"`
	Deprecated  bool                                     `raml:"deprecated"`
}

func MyHanderWithoutComment(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(MyStruct{})
}

// A route that use a handler fully commented
// @method PATCH
// @resource /myroute/{id}
// @route {string} id - The id of my route
// @query {bool} working - If set to `true`, everything works just fine
// @response {MyStruct}
// @example When everything works fine
//	/myroute/1?working=true
//	200: {
//		"value_1": "Hello world!",
// 		"value_2": 1,
// 		"value_3": true,
//	 	"value_4": {
//	 		"value_5": "2017-08-30T16:25:23.719Z",
//	 		"value_6": { },
// 		}
//	}
// @auth
// @deprecated
func MyHanderWithAllTheComments(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(MyStruct{})
}

// A route that use a handler partially commented
// @route {string} id - The id of my route
// @response {MyStruct}
func MyHanderWithFewComments(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(MyStruct{})
}

type MyStruct struct {
	Value1 string     `json:"value_1"`
	Value2 int        `json:"value_2"`
	Value3 bool       `json:"value_3"`
	Value4 *MyStruct2 `json:"value_4,omitempty"`
}

type MyStruct2 struct {
	Value5 []time.Time            `json:"value_5"`
	Value6 map[string]interface{} `json:"value_6"`
}
