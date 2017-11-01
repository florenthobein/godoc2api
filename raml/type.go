// RESTful API methods are operations that are performed on a resource
//
// Inspired by RAML 1.0 specs
// https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md/#raml-data-types
// and https://github.com/Jumpscale/go-raml/tree/master/raml

package raml

type Type struct {

	// Identifier of the type. (helper)
	Name string `yaml:"-"`

	// A default value for a type
	Default interface{} `yaml:"default,omitempty"`

	// Alias for the equivalent "type" property,
	// for compatibility with RAML 0.8.
	// Deprecated - API definitions should use the "type" property,
	// as the "schema" alias for that property name may be removed in a future RAML version.
	// The "type" property allows for XML and.
	Schema interface{} `yaml:"schema,omitempty"`

	// A base type which the current type extends,
	// or more generally a type expression.
	// A base type which the current type extends or just wraps.
	// The value of a type node MUST be either :
	//    a) the name of a user-defined type or
	//    b) the name of a built-in RAML data type (object, array, or one of the scalar types) or
	//    c) an inline type declaration.
	Type interface{} `yaml:"type,omitempty"`

	// An example of an instance of this type.
	// This can be used, e.g., by documentation generators to generate sample values for an object of this type.
	// Cannot be present if the examples property is present.
	// An example of an instance of this type that can be used,
	// for example, by documentation generators to generate sample values for an object of this type.
	// The "example" property MUST not be available when the "examples" property is already defined.
	Example interface{} `yaml:"example,omitempty"`

	// An object containing named examples of instances of this type.
	// This can be used, for example, by documentation generators
	// to generate sample values for an object of this type.
	// The "examples" property MUST not be available
	// when the "example" property is already defined.
	Examples map[string]interface{} `yaml:"examples,omitempty"`

	// An alternate, human-friendly name for the type
	DisplayName string `yaml:"displayName,omitempty"`

	// A substantial, human-friendly description of the type.
	// Its value is a string and MAY be formatted using markdown.
	Description string `yaml:"description,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having a key that begins
	// with "(" and ends with ")" where the text enclosed in parentheses is the annotation name,
	// and the value is an instance of that annotation.
	Annotations map[string]Annotation `yaml:",inline,omitempty"`

	// A map of additional, user-defined restrictions that will be inherited
	// and applied by any extending subtype.
	// https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md/#user-defined-facets
	// TODO
	//////// Facets map[string]string `yaml:"facets,omitempty"`

	// The capability to configure XML serialization of this type instance.
	// https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md/#xml-serialization-of-type-instances
	// TODO
	//////// XML XML `yaml:"xml,omitempty"`

	// An enumeration of all the possible values of instances of this type.
	// The value is an array containing representations of these possible values;
	// an instance of this type MUST be equal to one of these values.
	Enum []AnyType `yaml:"enum,flow,omitempty"`

	// Type extensions
	ObjectType ObjectType `yaml:",inline,omitempty"`
	ArrayType  ArrayType  `yaml:",inline,omitempty"`
	StringType StringType `yaml:",inline,omitempty"`
	NumberType NumberType `yaml:",inline,omitempty"`
	DateType   DateType   `yaml:",inline,omitempty"`
	FileType   FileType   `yaml:",inline,omitempty"`
}

type AnyType interface{}

type ObjectType struct {
	// The properties that instances of this type can or must have.
	Properties map[string]interface{} `yaml:"properties,omitempty"`

	// The minimum number of properties allowed for instances of this type.
	MinProperties int64 `yaml:"minProperties,omitempty"`

	// The maximum number of properties allowed for instances of this type.
	MaxProperties int64 `yaml:"maxProperties,omitempty"`

	// A Boolean that indicates if an object instance has additional properties.
	AdditionalProperties bool `yaml:"additionalProperties,omitempty"`

	// Determines the concrete type of an individual object at runtime when,
	// for example, payloads contain ambiguous types due to unions or inheritance.
	// The value must match the name of one of the declared properties of a type.
	// Unsupported practices are inline type declarations and using discriminator
	// with non-scalar properties.
	Discriminator string `yaml:"discriminator,omitempty"`

	// Identifies the declaring type. Requires including a discriminator facet in the type declaration.
	// A valid value is an actual value that might identify the type of an
	// individual object and is unique in the hierarchy of the type.
	// Inline type declarations are not supported.
	DiscriminatorValue string `yaml:"discriminatorValue,omitempty"`
}

type ArrayType struct {
	// Boolean value that indicates if items in the array MUST be unique.
	UniqueItems bool `yaml:"uniqueItems,omitempty"`

	// Indicates the type all items in the array are inherited from.
	// Can be a reference to an existing type or an inline type declaration.
	Items string `yaml:"items,omitempty"`

	// Minimum amount of items in array. Value MUST be equal to or greater than 0.
	MinItems int64 `yaml:"minItems,omitempty"`

	// Maximum amount of items in array. Value MUST be equal to or greater than 0.
	MaxItems int64 `yaml:"maxItems,omitempty"`
}

// Scalar types

type StringType struct {

	// Regular expression that this string SHOULD match.
	Pattern *string `yaml:"pattern,omitempty"`

	// Minimum length of the string. Value MUST be equal to or greater than 0.
	// ALSO used by File
	MinLength *int `yaml:"minLength,omitempty"`

	// Maximum length of the string. Value MUST be equal to or greater than 0.
	// ALSO used by File
	MaxLength *int `yaml:"maxLength,omitempty"`
}

// Number or Interger
type NumberType struct {

	// The minimum value of the parameter.
	Minimum *float64 `yaml:"minimum,omitempty"`

	// The maximum value of the parameter.
	Maximum *float64 `yaml:"maximum,omitempty"`

	// The format of the value. The value MUST be one of the following:
	// int32, int64, int, long, float, double, int16, int8
	// ALSO used by Date
	Format string `yaml:"format,omitempty"`

	// A numeric instance is valid against "multipleOf" if the result
	// of dividing the instance by this keyword's value is an integer.
	MultipleOf *int64 `yaml:"multipleOf,omitempty"`
}

type DateType struct {
	// The "full-date" notation of RFC3339, namely yyyy-mm-dd.
	// Does not support time or time zone-offset notation.
	DateOnly string `yaml:"date-only,omitempty"`

	// The "partial-time" notation of RFC3339, namely hh:mm:ss[.ff...].
	// Does not support date or time zone-offset notation.
	TimeOnly string `yaml:"time-only,omitempty"`

	// Combined date-only and time-only with a separator of "T",
	// namely yyyy-mm-ddThh:mm:ss[.ff...]. Does not support a time zone offset.
	DatetimeOnly string `yaml:"datetime-only,omitempty"`

	// A timestamp in one of the following formats:
	// if the format is omitted or set to rfc3339, uses the "date-time" notation of RFC3339;
	// if format is set to rfc2616, uses the format defined in RFC2616.
	// The additional facet `format` MUST be available only when the type
	// equals datetime, and the value MUST be either rfc3339 or rfc2616.
	// Any other values are invalid.
	Datetime string `yaml:"datetime,omitempty"`
}

type FileType struct {
	// A list of valid content-type strings for the file.
	// The file type MUST be a valid value.
	FileTypes []string `yaml:"fileTypes,omitempty"`
}

// Useful types
type HTTPCode uint
