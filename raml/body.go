// Body
//
// Inspired by RAML 1.0 specs
// https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md/#using-xml-and-json-schema

package raml

// Either directly a Type
// or under the facet `application/json`
// or under the facet `text/xml`
type Body struct {
	Type Type  `yaml:",inline,omitempty"`
	JSON *Type `yaml:"application/json,omitempty"`
	XML  *Type `yaml:"text/xml,omitempty"`
}
