package doc2raml

import (
	"reflect"

	"github.com/cometapp/midgar/doc2raml/raml"
)

type Trait struct{}

func DefineTrait(name string, kind reflect.Kind) {
	// Store the keyword
	reserveKeyword(name, KEYWORD_TYPE_TRAIT)
}

func (t *Trait) fillToRAML(index *map[string]raml.Trait) error {
	return nil
}
