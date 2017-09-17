package doc2raml

import "github.com/cometapp/midgar/doc2raml/raml"

const (
	SECURITY_OAUTH_1 = iota
	SECURITY_OAUTH_2
	SECURITY_BASIC_AUTHENTICATION
	SECURITY_DIGEST_AUTHENTICATION
	SECURITY_PASS_THROUGH
	SECURITY_X_CUSTOM
)

type Security struct{}

func DefineSecurity(name string, t uint) {
	// Store the keyword
	reserveKeyword(name, KEYWORD_TYPE_SECURITY)
}

func (s *Security) fillToRAML(index *map[string]raml.SecurityScheme) error {
	return nil
}
