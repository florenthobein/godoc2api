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

var securities_types_names_default = map[int]string{
	SECURITY_OAUTH_1:               "OAuth 1.0",
	SECURITY_OAUTH_2:               "OAuth 2.0",
	SECURITY_BASIC_AUTHENTICATION:  "Basic Authentication",
	SECURITY_DIGEST_AUTHENTICATION: "Digest Authentication",
	SECURITY_PASS_THROUGH:          "Pass Through",
	SECURITY_X_CUSTOM:              "x-custom",
}

// Registry of securities
var index_securities map[string]Security

// Security scheme, mirror of the RAML equivalent
type Security struct {
	Type            int
	TypeName        string // optional, only for SECURITY_X_CUSTOM, should start with `x-`
	Description     string
	Headers         map[string]Parameter
	QueryParameters map[string]Parameter
	QueryString     Parameter
	// Responses       []Response
	// Settings
}

// Configure a new security scheme.
// All the routes that declare the tag `tag_name` will be considered
// secured by this scheme.
func DefineSecurity(tag_name string, s Security) {
	// Store the keyword
	reserveTag(tag_name, _TAG_TYPE_SECURITY)
	// Store in the index
	if index_securities == nil {
		index_securities = map[string]Security{}
	}
	index_securities[tag_name] = s
}

// func (s *Security) fillToRAML(index *map[string]raml.SecurityScheme) error {
func securitiesToRAML(index *map[string]raml.SecurityScheme) error {
	if index == nil {
		return nil
	}

	for key, s := range index_securities {
		type_name := securities_types_names_default[s.Type]
		if s.TypeName != "" && len(s.TypeName) > 2 && s.TypeName[0] == 'x' && s.TypeName[1] == '-' {
			type_name = s.TypeName
		}
		description, has_description := raml.SecuritySchemeDescription{}, false
		if s.Headers != nil {
			has_description = true
			description.Headers = map[string]raml.Type{}
			for name, h := range s.Headers {
				description.Headers[name] = raml.Type{
					Description: h.Description,
					Example:     h.Example,
				}
			}
		}
		if s.QueryParameters != nil {
			has_description = true
			description.QueryParameters = map[string]raml.Type{}
			for name, qp := range s.QueryParameters {
				description.QueryParameters[name] = raml.Type{
					Description: qp.Description,
					Example:     qp.Example,
				}
			}
		}
		// todo query string
		ss := raml.SecurityScheme{
			Type:        type_name,
			Description: s.Description,
		}
		if has_description {
			ss.DescribedBy = description
		}
		(*index)[key] = ss
	}

	return nil
}
