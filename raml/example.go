// It is highly RECOMMENDED that API documentation include a rich selection of examples.
// RAML supports either the definition of multiple examples or a single one for any given
// instance of a type declaration.
//
// Inspired by RAML 1.0 specs
// https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md/#defining-examples-in-raml

package raml

type Example struct {

	// Identifier for the example. (helper)
	Name string `yaml:"-"`

	// An alternate, human-friendly name for the example.
	// If the example is part of an examples node, the default
	// value is the unique identifier that is defined for this example.
	DisplayName string `yaml:"displayName,omitempty"`

	// A substantial, human-friendly description for an example.
	// Its value is a string and MAY be formatted using markdown.
	Description string `yaml:"description,omitempty"`

	// The actual example of a type instance.
	Value interface{} `yaml:",omitempty"`

	// Validates the example against any type declaration (the default),
	// or not. Set this to false avoid validation.
	Strict bool `yaml:"strict"`
}
