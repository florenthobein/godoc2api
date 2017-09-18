package godoc2api

import (
	"reflect"

	"github.com/florenthobein/godoc2api/raml"
)

// Annotation type, mirror of the RAML equivalent
type Annotation struct{}

// Configure a new annotation type.
// All the routes that declare the tag `tag_name` will receive
// this annotation.
func DefineAnnotation(tag_name string, kind reflect.Kind) {
	// Store the keyword
	reserveTag(tag_name, _TAG_TYPE_ANNOTATION)

	// todo
}

func (a *Annotation) fillToRAML(index *map[string]raml.AnnotationType) error {
	return nil
}
