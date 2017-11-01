// The root section of the RAML document
//
// Inspired by RAML 1.0 specs
// https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md
// and https://github.com/Jumpscale/go-raml/tree/master/raml

package raml

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

const RAML_VERSION = "#%RAML 1.0"

// The root section of the RAML document describes the basic information about an API, such as its title and version.
// The root section also defines assets used elsewhere in the RAML document, such as types and traits.
type Root struct {

	// A short, plain-text label for the API. Its value is a string.
	Title string `yaml:"title"`

	// A substantial, human-friendly description of the API.
	// Its value is a string and MAY be formatted using markdown.
	Description string `yaml:"description,omitempty"`

	// The version of the API, for example "v1". Its value is a string.
	Version string `yaml:"version,omitempty"`

	// A URI that serves as the base for URIs of all resources.
	// Often used as the base of the URL of each resource containing the location of the API.
	// Can be a template URI.
	// The OPTIONAL baseUri property specifies a URI as an identifier for the API as a whole,
	// and MAY be used the specify the URL at which the API is served (its service endpoint),
	// and which forms the base of the URLs of each of its resources.
	// The baseUri property's value is a string that MUST conform to the URI specification RFC2396 or a Template URI.
	BaseURI string `yaml:"baseUri,omitempty"`

	// Named parameters used in the baseUri (template).
	// Any other URI template variables (than version) appearing in the baseUri MAY be described
	// explicitly within a baseUriParameters node at the root of the API definition.
	// The baseUriParameters node has the same structure and semantics as
	// the uriParameters node on a resource node, except that it specifies parameters
	// in the base URI rather than the relative URI of a resource.
	BaseURIParameters map[string]Parameter `yaml:"baseUriParameters,omitempty"`

	// The protocols supported by the API.
	// The OPTIONAL protocols property specifies the protocols that an API supports.
	// If the protocols property is not explicitly specified, one or more protocols
	// included in the baseUri property is used;
	// if the protocols property is explicitly specified,
	// the property specification overrides any protocol included in the baseUri property.
	// The protocols property MUST be a non-empty array of strings, of values HTTP and/or HTTPS, and is case-insensitive.
	Protocols []string `yaml:"protocols,omitempty"`

	// The default media types to use for request and response bodies (payloads),
	// for example "application/json".
	// Specifying the OPTIONAL mediaType property sets the default for return by API
	// requests having a body and for the expected responses. You do not need to specify the media type within every body definition.
	// The value of the mediaType property MUST be a sequence of
	// media type strings or a single media type string.
	// The media type applies to requests having a body,
	// the expected responses, and examples using the same sequence of media type strings.
	// Each value needs to conform to the media type specification in RFC6838.
	MediaType string `yaml:"mediaType,omitempty"`

	// Additional overall documentation for the API.
	// The API definition can include a variety of documents that serve as a
	// user guides and reference documentation for the API. Such documents can
	// clarify how the API works or provide business context.
	// The value of the documentation node is a sequence of one or more documents.
	// Each document is a map that MUST have exactly two key-value pairs: title and content
	Documentation []map[string]string `yaml:"documentation,omitempty"`

	// An alias for the equivalent "types" property for compatibility with RAML 0.8.
	// Deprecated - API definitions should use the "types" property
	// because a future RAML version might remove the "schemas" alias for that property name.
	// The "types" property supports XML and JSON schemas.
	// TODO
	//////// Schemas []map[string]string `yaml:"-"`

	// Declarations of (data) types for use within the API.
	Types map[string]Type `yaml:"types,omitempty"`

	// Declarations of traits for use within the API.
	Traits map[string]Trait `yaml:"traits,omitempty"`

	// Declarations of resource types for use within the API.
	// TODO
	//////// ResourceTypes map[string]interface{} `yaml:"resourceTypes,omitempty"`

	// Declarations of annotation types for use by annotations.
	// The value of the annotationsType node is a map whose keys define annotation type names,
	// also referred to as annotations, and whose values are key-value pairs called
	// annotation type declarations.
	AnnotationTypes map[string]AnnotationType `yaml:"annotationTypes,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having a key that begins
	// with "(" and ends with ")" where the text enclosed in parentheses is the annotation name,
	// and the value is an instance of that annotation.
	// TODO
	//////// Annotations map[string]Annotation `yaml:",inline,omitempty"`

	// Declarations of security schemes for use within the API.
	SecuritySchemes map[string]SecurityScheme `yaml:"securitySchemes,omitempty"`

	// The security schemes that apply to every resource and method in the API.
	SecuredBy []string `yaml:"securedBy,flow,omitempty"`

	// Imported external libraries for use within the API.
	// TODO
	//////// Uses map[string]string `yaml:"uses,omitempty"`

	// The resources of the API, identified as relative URIs that begin with a slash (/).
	// A resource property is one that begins with the slash and is either
	// at the root of the API definition or a child of a resource property. For example, /users and /{groupId}.
	Resources map[string]Resource `yaml:",inline,omitempty"`
}

// Check the coherence of the structure
// according to the RAML 1.0 specs
func (root *Root) Check() (bool, []error) {

	// Check annotations
	if len(root.AnnotationTypes) != 0 {
		// TODO
	}

	// Check types
	if len(root.Types) != 0 {
		// TODO
	}

	// Check resources
	if len(root.Resources) != 0 {
		// TODO
	}

	return true, nil
}

// Sort URI strings by their weight (quantity of /)
// then sort alphabetically
// Ex:
// 		/a < /b
// 		/b < /a/c
type ByURI []string

func (s ByURI) Len() int {
	return len(s)
}
func (s ByURI) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByURI) Less(i, j int) bool {
	wi := strings.Count(s[i], "/")
	wj := strings.Count(s[j], "/")
	if wi != wj {
		return wi < wj
	}
	return strings.Compare(s[i], s[j]) < 0
}

// Transform a flat list of resources into a tree-shaped list
func (root *Root) PileResources() {

	if len(root.Resources) == 0 {
		return
	}

	// Create an index of resources ordered by URI
	index := make(map[string]*Resource)
	ordered := make([]string, len(root.Resources))
	i := 0
	for uri, _ := range root.Resources {
		ordered[i] = uri
		i++
	}
	sort.Sort(ByURI(ordered))
	for _, uri := range ordered {
		r := root.Resources[uri]
		index[uri] = &r
	}

	// Filter the URI parameters of a RAML resource with a string URI
	var filterURIParameters = func(r *Resource, filter string) {
		ps := regexp.MustCompile(`\{([^\}]+)\}`).FindStringSubmatch(filter)
		if len(ps) <= 1 {
			(*r).URIParameters = nil
			return
		}
		ps = ps[1:len(ps)]
		var filtered map[string]Type
		// Get only the ones that are in URIParameters
		// and in the filter
		for _, p := range ps {
			if v, ok := (*r).URIParameters[p]; ok {
				if filtered == nil {
					filtered = make(map[string]Type)
				}
				filtered[p] = v
			}
		}
		(*r).URIParameters = filtered
	}

	// Create links between resources
	for current, r := range index {
		broken := strings.Split(r.URI, "/")
		max := len(broken)
		for i := max - 1; i > 1; i-- {
			base := strings.Join(broken[:i], "/")
			end := "/" + strings.Join(broken[i:max], "/")
			if _, ok := index[base]; ok {

				// Save the child in the parent
				p := index[base]
				if p.NestedResources == nil {
					p.NestedResources = make(map[string]*Resource)
				}
				p.NestedResources[end] = r

				// Save the parent in the child
				r.Parent = p

				// Filter the URI parameters
				filterURIParameters(r, end)

				// Store in the index
				index[base] = p
				index[current] = r
				break
			}
		}
	}

	// Get the root resources
	root_resources := make(map[string]Resource)
	for k, r := range index {
		if r.Parent == nil {
			// Filter the URI parameters for the root resources
			filterURIParameters(r, r.URI)
			root_resources[k] = *r
		}
	}
	root.Resources = root_resources

	return
}

// Return a string description of the document
func (root *Root) String() string {

	// Marshal the RAML document
	b, err := yaml.Marshal(root)
	if err != nil {
		log.Print(err)
	}

	return fmt.Sprintf("%s\n---\n%s", RAML_VERSION, b)
}
