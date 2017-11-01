// Response
//
// Inspired by RAML 1.0 specs
// https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md/#responses

package raml

type Response struct {

	// Identifier of the response. (helper)
	HTTPCode HTTPCode `yaml:"-"`

	// A substantial, human-friendly description of a response.
	// Its value is a string and MAY be formatted using markdown.
	Description string `yaml:"description,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having a key that begins
	// with "(" and ends with ")" where the text enclosed in parentheses is the annotation name,
	// and the value is an instance of that annotation.
	Annotations map[string]Annotation `yaml:",inline,omitempty"`

	// An API's methods may support custom header values in responses
	// Detailed information about any response headers returned by this method
	// TODO
	//////// Headers map[string]Header `yaml:"headers,omitempty"`

	// The body of the response
	Body Body `yaml:"body"`
}
