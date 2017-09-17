// Annotations provide a mechanism to extend the API specification with metadata
// beyond the metadata already defined in this RAML 1.0 specification.
//
// Inspired by RAML 1.0 specs
// https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md/#annotations

package raml

type AnnotationType struct {

	// Identifier for the annotation type. (helper)
	Name string `yaml:"-"`

	// An alternate, human-friendly method name in the context of the resource.
	// If the displayName property is not defined for a method,
	// documentation tools SHOULD refer to the resource by its property key,
	// which acts as the method name.
	DisplayName string `yaml:"displayName,omitempty"`

	// A longer, human-friendly description of the method in the context of the resource.
	// Its value is a string and MAY be formatted using markdown.
	Description string `yaml:"description,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having a key that begins
	// with "(" and ends with ")" where the text enclosed in parentheses is the annotation name,
	// and the value is an instance of that annotation.
	Annotations map[string]Annotation `yaml:",inline,omitempty"`

	// Detailed information about any query parameters needed by this method.
	// Mutually exclusive with queryString.
	// The value of the queryParameters node is a properties declaration object,
	// as is the value of the properties object of a type declaration
	AllowedTargets []TargetLocation `yaml:"allowedTargets,flow,omitempty"`

	// Type
	Type Type `yaml:",inline,omitempty"`
}

type TargetLocation string

var AllowedTargetLocation []string = []string{
	"API",                    // The root of a RAML document
	"DocumentationItem",      // An item in the collection of items that is the value of the root-level documentation node
	"Resource",               // A resource (relative URI) node, root-level or nested
	"Method",                 // A method node
	"Response",               // A declaration of the responses node, whose key is an HTTP status code
	"RequestBody",            // The body node of a method
	"ResponseBody",           // The body node of a response
	"TypeDeclaration",        // A data type declaration (inline or in a global types collection), header declaration, query parameter declaration, URI parameter declaration, or a property within any of these declarations, where the type property can be used
	"Example",                // Either an example or examples node
	"ResourceType",           // A resource type node
	"Trait",                  // A trait node
	"SecurityScheme",         // A security scheme declaration
	"SecuritySchemeSettings", // The settings node of a security scheme declaration
	"AnnotationType",         // A declaration of the annotationTypes node, whose key is a name of an annotation type and whose value describes the annotation
	"Library",                // The root of a library
	"Overlay",                // The root of an overlay
	"Extension",              // The root of an extension
}
