// Test the library
package godoc2api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/florenthobein/godoc2api"
)

// Configuration before tests
func init() {
	godoc2api.LogLevel = godoc2api.LOG_DEBUG

	// SecuritySchemes
	godoc2api.DefineSecurity("auth", godoc2api.Security{
		Type:        godoc2api.SECURITY_X_CUSTOM,
		TypeName:    "x-bearer",
		Description: "Authenticate a user with her auth token in the header",
		Headers: map[string]godoc2api.Parameter{
			"Authorization": godoc2api.Parameter{
				Description: "The user auth token preceded by Bearer",
				Example:     "Bearer _token_",
			},
		},
	})

	// Traits
	godoc2api.DefineTrait("pagination", nil) // todo

	// Annotations
	godoc2api.DefineAnnotation("deprecated", nil) // todo

	// Types
	godoc2api.DefineType("MyStruct", MyStruct{})
	godoc2api.DefineType("MyStruct2", MyStruct2{})

	return
}

// Compare & teardown
func finalize(folder string, t *testing.T) {
	compare(folder, t)
	teardown(folder)
}

// Compare results with fixtures at the end of the test
func compare(folder string, t *testing.T) {
	// Get fixture
	fixture, err := ioutil.ReadFile("fixtures/" + folder + "/test_api_v1.raml")
	if err != nil {
		return
	}
	// Get the result
	result, err := ioutil.ReadFile(folder + "/test_api_v1.raml")
	if err != nil {
		return
	}
	if string(result) != string(fixture) {
		t.Errorf("unexpected result for %s", folder)
	}
}

// Teardown after tests
func teardown(folder string) {
	os.RemoveAll(folder)
}

func TestWithRouteDefinition(t *testing.T) {
	output_dir := "test1"
	defer finalize(output_dir, t)

	doc := godoc2api.Documentation{
		Title:       "Test API",
		Description: "API used for tests",
		Version:     "v1",
		URL:         "http://mywebsite/{version}",
	}

	// Add a full route definition, the handler is not commented
	err := doc.AddRoute(RouteDefinition{
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
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	err = doc.Save(output_dir)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

}

func TestWithHandler(t *testing.T) {
	output_dir := "test2"
	defer finalize(output_dir, t)

	doc := godoc2api.Documentation{
		Title:       "Test API",
		Description: "API used for tests",
		Version:     "v1",
		URL:         "http://mywebsite/{version}",
	}

	// Add a simple handler without additional comment
	err := doc.AddRoute(MyHanderWithAllTheComments)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	err = doc.Save(output_dir)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
}

func TestBalanced(t *testing.T) {
	output_dir := "test3"
	defer finalize(output_dir, t)

	doc := godoc2api.Documentation{
		Title:       "Test API",
		Description: "API used for tests",
		Version:     "v1",
		URL:         "http://mywebsite/{version}",
	}

	// A more balanced approach
	err := doc.AddRoute(RouteDefinition{
		Resource: "GET /myroute/{id}",
		Handler:  MyHanderWithFewComments,
		Auth:     true,
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	err = doc.Save(output_dir)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
}

type RouteDefinition struct {
	Method      string           `raml:"method"`
	Resource    string           `raml:"resource"`
	Description string           `raml:"description"`
	Handler     http.HandlerFunc `raml:"handler"`
	RouteParams [][]string       `raml:"routes"`
	QueryParams [][]string       `raml:"queries"`
	Body        string           `raml:"body"`
	Response    string           `raml:"response"`
	Examples    [][]string       `raml:"examples"`
	Auth        bool             `raml:"auth"`
	Deprecated  bool             `raml:"deprecated"`
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
