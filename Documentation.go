// Godoc2API

// todo
// Fix: combinable enums is not RAML 1.0 compliant
// Improvement: Handle array type definitions like: (string | Person)[] (https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md/#type-expressions)
// Improvement: Create files for types to include

package godoc2api

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/florenthobein/godoc2api/raml"
)

// Default settings used for rendering the RAML root document
const (
	_DEFAULT_TITLE      = "Your API"
	_DEFAULT_VERSION    = "v1"
	_DEFAULT_URL        = "http://localhost/{v1}"
	_DEFAULT_MEDIA_TYPE = "application/json"
	_DEFAULT_OUTPUT_DIR = "raml"
)

// The tag that is used to parse an evenutal input struct that describes a route
const _MAIN_TAG_NAME = "raml"

// Main documentation struct, used to render a complete RAML documentation.
type Documentation struct {
	Title       string
	Description string
	Version     string
	URL         string
	MediaType   string
	routes      map[string]Route
	types       map[string]Type
	traits      map[string]Trait
	securities  map[string]Security
	annotations map[string]Annotation
}

// Add a route to the documentation.
//
// Adding a route via a handler function
//
// A route can be registered throught its http handler function, in that case the function should be fully documented, including with the tag `@resource` and `@method` (cf Reserved tags https://godoc.org/github.com/florenthobein/godoc2api/#pkg-constants).
//
// Example:
//	// Define a http handler
//	// @resource GET /myroute
//	// @response {MyObject}
//	func MyHandler(http.ResponseWriter, *http.Request) { ... }
//
//	...
//
//	// Register the route
//	d := godoc2api.Documentation{}
//	d.AddRoute(MyHandler)
//
// Adding a route via a struct
//
// A route can also be registered via a struct which fields use the tag `raml`, matching the reserved tags (cf Reserved tags https://godoc.org/github.com/florenthobein/godoc2api/#pkg-constants).
//
// Example:
//	// Define a http handler
//	func MyHandler(http.ResponseWriter, *http.Request) { ... }
//
//	...
//
//	// Custom route definition struct
//	type MyRouteDefinition struct {
//		Resource string `raml:"resource"`
//		Resp     string `raml:"response"`
//	}
//
//	...
//
//	// Register the route
//	d := godoc2api.Documentation{}
//	d.AddRoute(MyRouteDefinition{ "GET /myroute", "MyObject" })
func (d *Documentation) AddRoute(user_route interface{}) error {

	r := Route{_documentation: d}

	// Read the comment
	c, extra, err := readComment(user_route)
	if err != nil {
		warn("%v (%v)", err, user_route)
		return err
	}

	// If `user_route` is a struct and tags are already defined inside,
	// fill the route with it
	if extra != nil {
		for k, v := range extra {
			err := r.addTag(k, v)
			if err != nil {
				warn("%v (%v)", err, user_route)
			}
		}
	}

	// Parse the comment to extract tags and add it to the route
	tags := parseComment(c)
	for keyword, values := range tags {
		for _, fields := range values {
			err := r.addTag(keyword, fields)
			if err != nil {
				warn("%v (%v)", err, user_route)
			}
		}
	}

	// Check if the route can be used
	if err := r.checkViability(); err != nil {
		warn("unusable route: %v (%v)", err, user_route)
		return err
	}

	// Check if the route is coherent
	if err := r.checkURIParameters(); err != nil {
		warn("uncoherent route: %v (%s)", err, r.Resource)
	}

	// Store the route
	if d.routes == nil {
		d.routes = make(map[string]Route)
	}
	d.routes[r.signature()] = r

	return nil
}

// Generate the documentation
func (d *Documentation) toString() (string, error) {
	// Fill the empty fields
	if d.Title == "" {
		d.Title = _DEFAULT_TITLE
	}
	if d.Version == "" {
		d.Version = _DEFAULT_VERSION
	}
	if d.URL == "" {
		d.URL = _DEFAULT_URL
	}
	if d.MediaType == "" {
		d.MediaType = _DEFAULT_MEDIA_TYPE
	}

	// Format the document to a RAML structure
	api, err := d.toRAML()
	if err != nil {
		problem(err.Error())
		return "", err
	}

	// Transform the RAML into a string
	s := api.String()

	return s, nil
}

// Print the documentation
func (d *Documentation) Print() error {
	s, err := d.toString()
	if err != nil {
		return err
	}
	log.Println(s)
	return nil
}

// Render the documentation into a RAML file in the designated directory.
func (d *Documentation) Save(dirname string) error {

	sep := string(filepath.Separator)

	if dirname == "" {
		warn("no output directory specified, rendering to default %s", _DEFAULT_OUTPUT_DIR)
		dirname = _DEFAULT_OUTPUT_DIR
	} else {
		dirname = strings.Trim(dirname, " "+sep)
	}

	s, err := d.toString()
	if err != nil {
		return err
	}

	// Get the filename
	filename := fmt.Sprintf(
		"%s_%s.raml",
		regexp.MustCompile(`[^0-9a-z]`).ReplaceAllString(strings.ToLower(d.Title), "_"),
		d.Version,
	)

	// Create the directory
	dirpath := fmt.Sprintf(".%s%s",
		sep,
		dirname,
	)
	os.Mkdir(dirpath, 0777)

	// Create the file
	filepath := fmt.Sprintf("%s%s%s",
		dirpath,
		sep,
		filename,
	)
	err = ioutil.WriteFile(filepath, []byte(s), 0644)
	if err != nil {
		problem(err.Error())
	}
	return err
}

// The Types to add in the RAML document when rendering
func (d *Documentation) addType(t Type) bool {
	if d.types == nil {
		d.types = make(map[string]Type)
	}
	if _, ok := d.types[string(t)]; ok {
		return false
	}
	d.types[string(t)] = t

	other_ts := extractTypes(string(t))
	for _, t := range other_ts {
		d.addType(t)
	}

	return true
}

// Transform the documentation into a RAML structure
func (d *Documentation) toRAML() (raml.Root, error) {
	api := raml.Root{
		Title:       d.Title,
		Description: d.Description,
		Version:     d.Version,
		BaseURI:     d.URL,
		MediaType:   d.MediaType,
	}

	// Create the resources
	if d.routes != nil {
		api.Resources = make(map[string]raml.Resource)
		for _, r := range d.routes {
			err := r.fillToRAML(&api.Resources)
			if err != nil {
				return api, fmt.Errorf("error while RAMLing resource %s: %v", r.Resource, err)
			}
		}
		// Pile the resources
		api.PileResources()
	}

	// Create the types
	if d.types != nil {
		api.Types = make(map[string]raml.Type)
		for _, t := range d.types {
			err := t.fillToRAML(&api.Types)
			if err != nil {
				return api, fmt.Errorf("error while RAMLing type %s: %v", t, err)
			}
		}
	}

	// Create the annotation types
	if d.annotations != nil {
		api.AnnotationTypes = make(map[string]raml.AnnotationType)
		for _, a := range d.annotations {
			err := a.fillToRAML(&api.AnnotationTypes)
			if err != nil {
				return api, fmt.Errorf("error while RAMLing annotation %s: %v", a, err)
			}
		}
	}

	// Create the traits
	if d.traits != nil {
		api.Traits = make(map[string]raml.Trait)
		for _, t := range d.traits {
			err := t.fillToRAML(&api.Traits)
			if err != nil {
				return api, fmt.Errorf("error while RAMLing trait %s: %v", t, err)
			}
		}
	}

	// Create the security schemes
	if d.securities != nil {
		api.SecuritySchemes = make(map[string]raml.SecurityScheme)
		for _, s := range d.securities {
			err := s.fillToRAML(&api.SecuritySchemes)
			if err != nil {
				return api, fmt.Errorf("error while RAMLing security %s: %v", s, err)
			}
		}
	}

	return api, nil
}
