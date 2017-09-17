package doc2raml

import (
	"reflect"

	"github.com/cometapp/midgar/doc2raml/raml"
)

type Annotation struct{}

func DefineAnnotation(name string, kind reflect.Kind) {
	// Store the keyword
	reserveKeyword(name, KEYWORD_TYPE_ANNOTATION)
}

func (a *Annotation) fillToRAML(index *map[string]raml.AnnotationType) error {
	return nil
}
