package examples

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/florenthobein/godoc2api"
)

// A structure to describe your routes
type RouteDefinition struct {
	Method  string           `raml:"method"`
	URI     string           `raml:"resource"`
	Handler http.HandlerFunc `raml:"handler"`
	Auth    bool             `raml:"auth"`
}

func Example_bookshop() {

	// Define your routes
	routes := []RouteDefinition{
		RouteDefinition{"POST", "/books", CreateBooks, true},
		RouteDefinition{"POST", "/books/new", CreateBooksDEPRECATED, true},
		RouteDefinition{"GET", "/books/{id}", GetBook, false},
		RouteDefinition{"POST", "/books/{id}", UpdateBook, true},
		RouteDefinition{"DELETE", "/books/{id}", DeleteBook, true},
	}

	// Configure your documentation
	godoc2api.DefineSecurity("auth", nil)    // todo
	godoc2api.DefineTrait("deprecated", nil) // todo
	godoc2api.DefineType("Book", Book{})
	godoc2api.DefineTypeRAML("uuid", "string", map[string]interface{}{"pattern": `[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}`})
	doc := godoc2api.Documentation{
		Title: "Book collection",
		URL:   "http://localhost:8080",
	}

	for _, r := range routes {
		// Create your route
		handleFunc(r.Method, r.URI, r.Handler)
		// Add them to your doc
		doc.AddRoute(r)
	}

	// Save your doc
	err := doc.Save("example_bookshop/")
	if err != nil {
		panic(err)
	}

	// Run your webserver
	if mux.server != nil {
		log.Println("Go to http://localhost:8080")
		http.ListenAndServe(":8080", mux.server)
	}
}

// Handlers
// Where your API documentation really is
//////////////////////////////////////////////////////////////////////////////

// Create a new book
// @body {Book} The book you want to create
// @response {Book} The created book
// @example Create a classic
//	{
//		"id": "ca761232-ed42-11ce-bacd-00aa0057b223",
//		"name": "Cyrano de bergerac",
//		"author": "Edmond Rostand",
//		"price": 10.3,
//		"stars": 2
//	}
//	200: {
//		"id": "ca761232-ed42-11ce-bacd-00aa0057b223",
//		"name": "Cyrano de bergerac",
//		"author": "Edmond Rostand",
//		"price": 10.3,
//		"stars": 2
//	}
func CreateBooks(rw http.ResponseWriter, r *http.Request) {
	// ...doing things...
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(Book{})
}

// Create a new book the old way
// @body {Book} - The book you want to create
// @response {Book} - The created book
// @deprecated
func CreateBooksDEPRECATED(rw http.ResponseWriter, r *http.Request) {
	// ...doing things...
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(Book{})
}

// Get a book
// @route {uuid} id - The identifier of the book
// @query {bool} with_metadata - If set to `true`, includes metadatas in the response
// @response {Book} - The book that you wanted
func GetBook(rw http.ResponseWriter, r *http.Request) {
	// ...doing things...
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(Book{})
}

// Update a book
// @route {uuid} id - The identifier of the book
// @body {Book} - The book data to update
// @response {Book} - The book updated
func UpdateBook(rw http.ResponseWriter, r *http.Request) {
	// ...doing things...
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(Book{})
}

// Delete a book
// @route {uuid} id - The identifier of the book
// @body {Book} - The book data to update
// @response {Book} - The book updated
func DeleteBook(rw http.ResponseWriter, r *http.Request) {
	// ...doing things...
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(nil)
}

// Model
//////////////////////////////////////////////////////////////////////////////

// Define a simple object
type Book struct {
	Id          string  `json:"id" ramlType:"uuid"`
	Name        string  `json:"name"`
	Author      string  `json:"author"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
	Stars       uint    `json:"stars"`
}

// Mini multiplexer
//////////////////////////////////////////////////////////////////////////////

type Mux struct {
	server *http.ServeMux
	routes map[string]map[string]http.HandlerFunc // map[uri]map[method]handler
}

var mux = Mux{}

func (m *Mux) serve(URI string) (fn http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if m == nil {
			http.NotFound(w, r)
			return
		}
		if _, ok := (*m).routes[URI]; !ok {
			http.NotFound(w, r)
			return
		}
		if _, ok := (*m).routes[URI][method]; !ok {
			http.NotFound(w, r)
			return
		}
		(*m).routes[URI][method](w, r)
	}
}

// Helper that allow defining a route with a combination of `method` & `URI`
func handleFunc(method string, URI string, fn func(http.ResponseWriter, *http.Request)) {
	if mux.routes == nil {
		mux.server = http.NewServeMux()
		mux.routes = map[string]map[string]http.HandlerFunc{}
	}
	res := regexp.MustCompile(`^(.+){(.+)}$`).FindStringSubmatch(URI)
	if len(res) > 1 {
		URI = res[1]
	}
	if mux.routes[URI] == nil {
		mux.routes[URI] = map[string]http.HandlerFunc{}
		mux.routes[URI][method] = fn
		mux.server.HandleFunc(URI, mux.serve(URI))
		return
	}
	mux.routes[URI][method] = fn
	return
}
