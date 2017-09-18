package godoc2api

import "github.com/florenthobein/godoc2api/raml"

// Security schemes
const (
	SECURITY_OAUTH_1 = iota
	SECURITY_OAUTH_2
	SECURITY_BASIC_AUTHENTICATION
	SECURITY_DIGEST_AUTHENTICATION
	SECURITY_PASS_THROUGH
	SECURITY_X_CUSTOM
)

// Security scheme, mirror of the RAML equivalent
type Security struct{}

// Configure a new security scheme.
// All the routes that declare the tag `tag_name` will be considered
// secured by this scheme.
func DefineSecurity(tag_name string, t interface{}) {
	// Store the keyword
	reserveTag(tag_name, _TAG_TYPE_SECURITY)
}

func (s *Security) fillToRAML(index *map[string]raml.SecurityScheme) error {
	return nil
}
