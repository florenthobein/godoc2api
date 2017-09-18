package godoc2api

import "github.com/florenthobein/godoc2api/raml"

// Trait, mirror of the RAML equivalent
type Trait struct{}

// Configure a new trait.
// All the routes that declare the tag tag_name will be considered
// using this trait.
func DefineTrait(tag_name string, t interface{}) {
	// Store the keyword
	reserveTag(tag_name, _TAG_TYPE_TRAIT)

	// todo
}

func (t *Trait) fillToRAML(index *map[string]raml.Trait) error {
	return nil
}
