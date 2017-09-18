// Resources are accessible endpoints of an API, described in the RAML document
//
// Inspired by RAML 1.0 specs
// https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md/#resources-and-nested-resources
// and https://github.com/Jumpscale/go-raml/tree/master/raml

package raml

// A resource is identified by its relative URI, which MUST begin with a slash ("/").
// Every node whose key begins with a slash, and is either at the root of the API definition
// or is the child node of a resource node, is such a resource node.
type Resource struct {

	// Identifier of the resource. (helper)
	URI string `yaml:"-"`

	// An alternate, human-friendly name for the resource.
	// If the displayName property is not defined for a resource,
	// documentation tools SHOULD refer to the resource by its property key
	// which acts as the resource name. For example, tools should refer to the relative URI /jobs.
	DisplayName string `yaml:"displayName,omitempty"`

	// A substantial, human-friendly description of a resource.
	// Its value is a string and MAY be formatted using markdown.
	Description string `yaml:"description,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having a key that begins
	// with "(" and ends with ")" where the text enclosed in parentheses is the annotation name,
	// and the value is an instance of that annotation.
	// TODO
	//////// Annotations map[string]Annotation `yaml:",inline,omitempty"`

	// Detailed information about any URI parameters of this resource.
	URIParameters map[string]Type `yaml:"uriParameters,omitempty"`

	// In a RESTful API, methods are operations that are performed on a
	// resource. A method MUST be one of the HTTP methods defined in the
	// HTTP version 1.1 specification [RFC2616] and its extension,
	// RFC5789 [RFC5789].
	Get     *Method `yaml:"get,omitempty"`
	Patch   *Method `yaml:"patch,omitempty"`
	Put     *Method `yaml:"put,omitempty"`
	Head    *Method `yaml:"head,omitempty"`
	Post    *Method `yaml:"post,omitempty"`
	Delete  *Method `yaml:"delete,omitempty"`
	Options *Method `yaml:"options,omitempty"`

	// A list of traits to apply to all methods declared (implicitly or explicitly) for this resource.
	// Individual methods can override this declaration.
	// TODO
	//////// Is []interface{} `yaml:"is,omitempty"`

	// The resource type that this resource inherits.
	// TODO
	//////// Type interface{} `yaml:"type,omitempty"`

	// The security schemes that apply to all methods declared (implicitly or explicitly) for this resource.
	// TODO
	//////// SecuredBy []interface{} `yaml:"securedBy,omitempty"`

	// A resource defined as a child node of another resource is called a nested resource.
	// The key of the child node is the URI of the nested resource relative to the
	// parent resource URI.
	NestedResources map[string]*Resource `yaml:",inline"`

	// If this is not nil, then this resource is a nested resource.
	Parent *Resource `yaml:"-"`

	// All methods of this resource. (helper)
	Methods []*Method `yaml:"-"`
}
