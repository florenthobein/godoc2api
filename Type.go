package godoc2api

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/fatih/structs"

	"github.com/florenthobein/godoc2api/raml"
)

// Registry of type definitions
var index_types map[string]TypeDefinition

// Regex to match maps
const _PARSE_MAP = `^map\[([^\]]+)\](.+)$`

// A Type is a global definition of a parameter's type.
// Every type that is not a scalar should generate
// a type definition.
type Type string

// A TypeDefinition defines the structure of a Type and is linked
// at the of the document
type TypeDefinition struct {
	name         string
	aliasFor     *TypeDefinition
	nameRAMLType string
	reflectType  *reflect.Type
	properties   map[string]interface{}
	mapKey       string
	mapValue     string
}

// Configure a new type definition.
//
// The objects used as inputs and outputs of the API
// have to match a defined type.
// This is mandatory, as go doesn't have a global type registry that would allow
// to create a variable out of its kind's name.
//
// Example
//
// This type definition
//	DefineType("MyStruct", MyStruct{})
// should be used to define the response type of
//	// My route description
//	// @resource GET /myroute
//	// @response {MyStruct}
//	func MyHandler(http.ResponseWriter, *http.Request) { ... }
func DefineType(name string, obj interface{}) {
	if index_types == nil {
		index_types = make(map[string]TypeDefinition)
	}
	ref := reflect.TypeOf(obj)
	true_name := ref.String()
	td := TypeDefinition{
		name:        name,
		reflectType: &ref,
	}
	index_types[name] = td
	index_types[true_name] = TypeDefinition{
		name:     true_name,
		aliasFor: &td,
	}
}

// Configure a new type definition already formated with standard RAML data types.
//
// See https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md/#raml-data-types
// for more informations on RAML types
func DefineTypeRAML(name, raml_type string, properties map[string]interface{}) {
	if index_types == nil {
		index_types = make(map[string]TypeDefinition)
	}
	index_types[name] = TypeDefinition{
		name:         name,
		nameRAMLType: raml_type,
		properties:   properties,
	}
}

func defineTypeMap(name, key, value string) *TypeDefinition {
	if index_types == nil {
		index_types = make(map[string]TypeDefinition)
	}
	td := TypeDefinition{
		name:     name,
		mapKey:   key,
		mapValue: value,
	}
	index_types[name] = td
	return &td
}

func isDefinedTypeRAML(name string) (t TypeDefinition, ok bool) {
	if index_types == nil {
		return
	}
	t, ok = index_types[name]
	ok = ok && t.nameRAMLType != ""
	return
}

// Configure a new alias for a type definition.
//
// Example
//
// This type definition
//	DefineType("MyStruct", MyStruct{})
//	DefineTypeAlias("super_struct", "MyStruct")
// can be used to define the response type of
//	// My route description
//	// @resource GET /myroute
//	// @response {super_struct}
//	func (http.ResponseWriter, *http.Request) { ... }
//
// Limitation
//
// Why not directy using `DefineType("super_struct", MyStruct{})` then?
// Because if an other struct contains a field of the type `MyStruct`,
// the library will not be able to find a reference of the object.
//	type MyStruct2 []MyStruct // this requires `MyStruct` and `MyStruct2` to be defined
func DefineTypeAlias(alias, name string) {
	if index_types == nil {
		index_types = make(map[string]TypeDefinition)
	}
	if _, ok := index_types[name]; !ok {
		warn("can't define alias `%s`: type `%s` doesn't exist", alias, name)
		return
	}
	td := index_types[name]
	index_types[alias] = TypeDefinition{
		name:     alias,
		aliasFor: &td,
	}
}

func isTypeAlias(alias string) (res string, ok bool) {
	if index_types != nil {
		if td, exists := index_types[alias]; exists {
			ok = td.aliasFor != nil
			if ok {
				res = td.aliasFor.name
			}
		}
	}
	return
}

func isDefinedType(name string) (t TypeDefinition, ok bool) {
	if index_types == nil {
		return
	}
	t, ok = index_types[name]
	return
}

func formatMapName(key, val string) string {
	key = strings.Replace(key, "[]", "Array", -1)
	val = strings.Replace(strings.Replace(val,
		"[]", "Array", -1),
		" ", "", -1)
	_, k, _, _ := formatType(key)
	_, v, _, _ := formatType(val)
	// return fmt.Sprintf("map_%s_%s", key, val)
	return fmt.Sprintf("map_%s_%s", string(k), string(v))
}

// Given a type, format it to the RAML Data TypeDefinition format
// and add specify if it should be registired in Types.
// Examples:
// 		interface{}					=> any		any											nil
// 		string							=> scalar	string									nil
// 		*time.Time					=> scalar	datetime								nil
// 		int16								=> scalar	number									nil
// 		map[string]MyObject	=> object	map_string_boolean			[]Type{"MyObject"}
// 		MyObject						=> object	MyObject								[]Type{"MyObject"}
// 		[]MyObject					=> array	MyObject[]							[]Type{"MyObject"}
// 		map[string][]bool		=> array	map_string_Arrayboolean	nil
func formatType(name string) (global string, precise Type, register_types []Type, err error) {

	// Just to return a correct name
	strs := regexp.MustCompile(` | `).Split(name, -1)
	if len(strs) > 1 {
		register_types = []Type{}
		for i, str := range strs {
			if str == "|" {
				continue
			}
			_, v, other_ts, _ := formatType(str)
			for _, other_t := range other_ts {
				register_types = append(register_types, other_t)
			}
			strs[i] = string(v)
		}
		precise = Type(strings.Join(strs, " "))
		return
	}

	if name == "" {
		return "nil", Type("nil"), nil, nil
	}

	if name == "interface{}" {
		return "any", Type("any"), nil, nil
	}

	// If reserved
	if td, ok := isDefinedTypeRAML(name); ok {
		t := Type(td.name)
		return td.nameRAMLType, t, []Type{t}, nil
	}

	// If alias
	if alias, ok := isTypeAlias(name); ok {
		if name == alias {
			return "", "", nil, fmt.Errorf("loop alias for %s", name)
		}
		global, precise, register_types, err = formatType(alias)
		return
	}

	// If it's a pointer
	if len(name) > 1 && name[0:1] == "*" {
		global, precise, register_types, err = formatType(name[1:])
		return
	}

	// If it's a slice
	if len(name) > 2 && name[0:2] == "[]" {
		global = "array"
		_, precise, register_types, err = formatType(name[2:])
		precise = Type(string(precise) + "[]")
		return
	}

	// If it's a map
	if res := regexp.MustCompile(_PARSE_MAP).FindStringSubmatch(name); len(res) > 2 {
		register_types = []Type{Type(name)}
		_, k, other_ts, _ := formatType(res[1])
		for _, other_t := range other_ts {
			register_types = append(register_types, other_t)
		}
		_, v, other_ts, _ := formatType(res[2])
		for _, other_t := range other_ts {
			register_types = append(register_types, other_t)
		}
		name = formatMapName(res[1], res[2])
		defineTypeMap(name, string(k), string(v))
		debug("creation of type map %s", name)
		return "object", Type(name), register_types, nil
	}

	// If it's a scalar
	if is_scalar, name := isScalar(name); is_scalar {
		return "scalar", Type(name), nil, nil
	}

	// If it's a nil
	if name == "nil" {
		return "nil", Type("nil"), nil, nil
	}

	t := Type(name)
	return "object", t, []Type{t}, nil
}

func isScalar(name string) (bool, string) {
	switch name {
	case "time.Time":
		return true, "datetime"
	case "bool":
		return true, "boolean"
	case "int", "int16", "int8", "int32", "int64",
		"uint", "uint16", "uint8", "uint32", "uint64":
		return true, "integer"
	case "float32", "float64":
		return true, "number"
	case "string":
		return true, "string"
	case "multipart.FileHeader":
		return true, "file"
	}
	return false, name
}

func (t *Type) fillToRAML(types *map[string]raml.Type) error {
	// Check the index
	if index_types == nil {
		return fmt.Errorf("no index type")
	}

	name := string(*t)

	// Check multiple types
	others := strings.Split(name, " | ")
	if len(others) > 1 {
		for _, other := range others {
			t := Type(other)
			t.fillToRAML(types)
		}
		return nil
	}
	name = strings.Trim(others[0], " ")

	// Check the alias
	alias := ""
	_ = alias
	if val, exists := isTypeAlias(name); exists {
		alias = name
		name = val
	}

	// Check map
	if res := regexp.MustCompile(_PARSE_MAP).FindStringSubmatch(name); len(res) > 2 {
		name = formatMapName(res[1], res[2])
	}

	// Check name of module
	// TODO tricky
	// if res := strings.Split(name, "."); len(res) > 1 {
	// 	name = res[1]
	// }

	// Check the type definition
	td, ok := index_types[name]
	if !ok {
		return fmt.Errorf("type `%s` not found", name)
	}

	// raml_type, others := td.toRAML()
	raml_type, _ := td.toRAML()
	(*types)[td.name] = raml_type
	// for _, other := range others {
	// 	t := Type(other)
	// 	t.fillToRAML(types)
	// }
	return nil
}

func (td *TypeDefinition) toRAML() (raml.Type, []string) {
	// Other types to generate
	others := []string{}

	description := ""

	// Check if this Type is defined as a RAML Type
	switch td.nameRAMLType {
	case "string":
		st := raml.StringType{}
		if v, ok := td.properties["description"]; ok {
			v_typed := v.(string)
			description = v_typed
		}
		if v, ok := td.properties["pattern"]; ok {
			v_typed := v.(string)
			st.Pattern = &v_typed
		}
		if v, ok := td.properties["minLength"]; ok {
			v_typed := v.(int)
			st.MinLength = &v_typed
		}
		if v, ok := td.properties["maxLength"]; ok {
			v_typed := v.(int)
			st.MaxLength = &v_typed
		}
		if v, ok := td.properties["length"]; ok {
			v_typed := v.(int)
			st.MinLength = &v_typed
			st.MaxLength = &v_typed
		}
		return raml.Type{Type: td.nameRAMLType, Description: description, StringType: st}, others
	}

	var mapToType = func(k, v string) raml.Type {
		reg := ""
		_, tk, _, err := formatType(k)
		if err != nil {
			warn(err.Error())
			return raml.Type{}
		}
		if t, ok := isDefinedTypeRAML(string(tk)); ok {
			tk = Type(t.nameRAMLType)
		}
		switch string(tk) {
		case "integer", "number":
			reg = "/^[0-9]+$/"
			break
		case "string":
			reg = "/^.*$/"
			break
		default:
			warn("unknown mapToType %s", k)
		}
		_, tv, _, err := formatType(v)
		if err != nil {
			warn(err.Error())
			return raml.Type{}
		}
		return raml.Type{
			Type: "object",
			ObjectType: raml.ObjectType{
				Properties: map[string]interface{}{
					reg: tv,
				},
				AdditionalProperties: true,
			},
		}
	}

	// If map
	if td.mapKey != "" && td.mapValue != "" {
		return mapToType(td.mapKey, td.mapValue), []string{}
	}

	// If not a raml type by default
	if td.reflectType == nil {
		warn("wrong type %v", td)
		return raml.Type{}, others
	}

	// Create a new obj
	v := reflect.New(*td.reflectType).Elem()
	instance := v.Interface()

	// If not a struct
	if !structs.IsStruct(instance) {
		name := ""
		typeof := reflect.TypeOf(instance)
		if v.Kind().String() == "map" {
			return mapToType(typeof.Key().String(), strings.Replace(typeof.Elem().String(), " ", "", -1)), others
		} else if v.Kind().String() == "slice" {
			name = "[]" + strings.Replace(typeof.Elem().String(), " ", "", -1)
		} else {
			warn("[%s] %s (%s)", v.Kind().String(), (*td.reflectType).String(), reflect.TypeOf(instance).Elem().String())
			return raml.Type{}, others
		}

		_, precise, _, err := formatType(name)
		if err != nil {
			warn(err.Error())
			return raml.Type{}, others
		}

		return raml.Type{Type: precise}, others
	}

	// Read the struct
	properties := map[string]interface{}{}
	alt_tag_name := "json"
	main_tag_type_name := "ramlType"
	s := structs.New(instance)
	fs := s.Fields()
	for _, f := range fs {
		value := f.Tag(_MAIN_TAG_NAME)
		if value == "" {
			value = f.Tag(alt_tag_name)
		}
		if value == "" || value == "-" {
			continue
		}
		v := strings.Split(value, ",")
		if len(v) < 1 || v[0] == "" {
			continue
		}
		name := v[0]

		// Check the kind
		type_name := strings.Replace(reflect.TypeOf(f.Value()).String(), " ", "", -1)
		if value := f.Tag(main_tag_type_name); value != "" {
			type_name = value
		}
		_, precise, _, err := formatType(type_name)
		if err != nil {
			warn(err.Error())
			continue
		}

		// Optional
		optional := strings.Contains(value, ",omitempty") || (len(type_name) > 1 && type_name[0:1] == "*" && strings.Contains(value, ",omitempty"))
		if optional {
			name = name + "?"
		}

		properties[name] = precise
	}

	return raml.Type{
		Type: "object",
		ObjectType: raml.ObjectType{
			Properties:           properties,
			AdditionalProperties: false,
		},
	}, others
}

func extractTypes(name string) (ts []Type) {
	var td TypeDefinition
	var ok bool
	if td, ok = isDefinedType(name); !ok || td.nameRAMLType != "" {
		return
	}

	if td.aliasFor != nil {
		td = *td.aliasFor
	}

	if td.reflectType == nil {
		return
	}

	// Create a new obj
	v := reflect.New(*td.reflectType).Elem()
	instance := v.Interface()

	var process = func(rt reflect.Type, kind string) []Type {
		res := []Type{}
		item := ""
		if kind == "map" {
			item = formatMapName(rt.Key().String(), strings.Replace(rt.Elem().String(), " ", "", -1))
		} else if kind == "slice" {
			item = strings.Replace(rt.Elem().String(), " ", "", -1)
		} else {
			item = strings.Replace(rt.String(), " ", "", -1)
		}
		_, _, register_types, _ := formatType(item)
		for _, new_t := range register_types {
			res = append(res, new_t)
		}
		return res
	}

	// If not a struct
	if !structs.IsStruct(instance) {
		return process(reflect.TypeOf(instance), v.Kind().String())
	}

	// Read the struct
	ts = []Type{}
	s := structs.New(instance)
	fs := s.Fields()
	for _, f := range fs {
		new_ts := process(reflect.TypeOf(f.Value()), v.Kind().String())
		for _, new_t := range new_ts {
			ts = append(ts, new_t)
		}
	}

	return
}
