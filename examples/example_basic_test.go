package examples

import (
	"fmt"
	"net/http"

	"github.com/florenthobein/godoc2api"
)

func Example_basic() {

	// Define your normal route
	http.HandleFunc("/myroute", myHander)

	// Define your documentation and save it
	doc := godoc2api.Documentation{URL: "http://localhost:8080"}
	doc.AddRoute(myHander)
	doc.Save("example_basic/")

	// Run your webserver
	http.ListenAndServe(":8080", nil)
}

// A simple route
// @resource GET /myroute
func myHander(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Hello world!")
}
