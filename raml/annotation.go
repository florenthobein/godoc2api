// Annotations provide a mechanism to extend the API specification with metadata
// beyond the metadata already defined in this RAML 1.0 specification.
//
// Inspired by RAML 1.0 specs
// https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md/#annotations

package raml

// To be applied in an API specification, the annotation MUST be declared in an annotation type.
type Annotation struct {

	// Identifier for the annotation. (helper)
	Name string `yaml:"-"`

	// Type
	Type
}

type Annotations map[string]Annotation
