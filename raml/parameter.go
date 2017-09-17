// Parameters

package raml

type Parameter struct {

	// Identifier for the parameter. (helper)
	Name string `yaml:"-"`

	// A friendly name used only for display or documentation purposes.
	// If displayName is not specified, it defaults to the property's key.
	DisplayName string `yaml:"displayName,omitempty"`

	// The intended use or meaning of the parameter.
	Description string `yaml:"description,omitempty"`
}
