// RESTful API methods are operations that are performed on a resource
//
// Inspired by RAML 1.0 specs
// https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md/#methods
// and https://github.com/Jumpscale/go-raml/tree/master/raml

package raml

// These correspond to the HTTP methods defined in the HTTP version 1.1 specification
// RFC2616 and its extension, RFC5789.
type Method struct {

	// Identifier for the method. (helper)
	Name string `yaml:"-"`

	// An alternate, human-friendly method name in the context of the resource.
	// If the displayName property is not defined for a method,
	// documentation tools SHOULD refer to the resource by its property key,
	// which acts as the method name.
	DisplayName string `yaml:"displayName,omitempty"`

	// A longer, human-friendly description of the method in the context of the resource.
	// Its value is a string and MAY be formatted using markdown.
	Description string `yaml:"description,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having a key that begins
	// with "(" and ends with ")" where the text enclosed in parentheses is the annotation name,
	// and the value is an instance of that annotation.
	Annotations map[string]Annotation `yaml:",inline,omitempty"`

	// Detailed information about any query parameters needed by this method.
	// Mutually exclusive with queryString.
	// The value of the queryParameters node is a properties declaration object,
	// as is the value of the properties object of a type declaration
	QueryParameters map[string]Type `yaml:"queryParameters,omitempty"`

	// Detailed information about any request headers needed by this method.
	// TODO
	//////// Headers map[HTTPHeader]Header `yaml:"headers,omitempty"`

	// The query string needed by this method.
	// Mutually exclusive with queryParameters.
	// TODO
	//////// QueryString map[string]Type `yaml:"queryString,omitempty"`

	// Information about the expected responses to a request.
	// Responses MUST be a map of one or more HTTP status codes, where each
	// status code itself is a map that describes that status code.
	Responses map[HTTPCode]Response `yaml:"responses,omitempty"`

	// A request body that the method admits.
	Body *Body `yaml:"body,omitempty"`

	// Explicitly specify the protocol(s) used to invoke a method,
	// thereby overriding the protocols set elsewhere,
	// for example in the baseUri or the root-level protocols property.
	// TODO
	//////// Protocols []string `yaml:"protocols,omitempty"`

	// A list of the traits to apply to this method.
	// TODO
	//////// Is []interface{} `yaml:"is,omitempty"`

	// The security schemes that apply to this method.
	// TODO
	//////// SecuredBy []interface{} `yaml:"securedBy,omitempty"`
}
