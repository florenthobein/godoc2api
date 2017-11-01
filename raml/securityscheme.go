// Security schemes
//
// Inspired by RAML 1.0 specs
// https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md/#security-schemes
// and https://github.com/Jumpscale/go-raml/tree/master/raml

package raml

type SecuritySchemeDescription struct {
	// Optional array of Headers, documenting the possible headers that could be accepted.
	Headers map[string]Type `yaml:"headers,omitempty"`

	// Query parameters, used by the schema to authorize the request.
	// Mutually exclusive with queryString.
	QueryParameters map[string]Type `yaml:"queryParameters,omitempty"`

	// The query string used by the schema to authorize the request.
	// Mutually exclusive with queryParameters.
	QueryString Type `yaml:"queryString,omitempty"`

	// An optional array of responses, representing the possible responses that could be sent.
	Responses map[HTTPCode]Response `yaml:"responses,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having a key that begins
	// with "(" and ends with ")" where the text enclosed in parentheses is the annotation name,
	// and the value is an instance of that annotation.
	Annotations map[string]Annotation `yaml:",inline,omitempty"`
}

type SecurityScheme struct {

	// Identifier for the security scheme. (helper)
	Name string `yaml:"-"`

	// An alternate, human-friendly name for the security scheme.
	DisplayName string `yaml:"displayName,omitempty"`

	// Specifies the API security mechanisms. One API-supported authentication method is allowed.
	// The value MUST be one of the following methods:
	// OAuth 1.0, OAuth 2.0, Basic Authentication, Digest Authentication, Pass Through, x-<other>
	Type string `yaml:"type"`

	// Information that MAY be used to describe a security scheme.
	// Its value is a string and MAY be formatted using markdown.
	Description string `yaml:"description,omitempty"`

	// A description of the following security-related request components determined by the scheme:
	// the headers, query parameters, or responses. As a best practice, even for standard
	// security schemes, API designers SHOULD describe these nodes of security schemes.
	// Including the security scheme description completes the API documentation.
	DescribedBy SecuritySchemeDescription `yaml:"describedBy,omitempty"`

	// The settings attribute MAY be used to provide security scheme-specific information.
	Settings string `yaml:"settings,omitempty"`
}
