// The parser is responsible for analysing a comment and extracting meaningful keywords
// It also structures the keyword fields

package godoc2api

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Regular expressions used for parsing the comments
const (
	_PARSE_METHODS         = `(GET|HEAD|POST|PUT|DELETE|PATCH|OPTIONS)`
	_PARSE_RESOURCE        = `^(?:` + _PARSE_METHODS + ` )?(/.+)$`
	_PARSE_TAG             = `^(?://| ?\*) @(\w+)(?:[ 	]+(.+))?$`
	_PARSE_TAGBLOCK        = `^(?://| ?\*)(?:[ 	]+(.+))?$`
	_PARSE_LINE            = `^\{\(?([^\)]+)\)?\}(?:[ 	]+\[?([\w\=]+)?\]?(?:[ 	\-]+(?:\-[ 	]+)?(.+))?)?$`
	_PARSE_TYPE            = `^([\w \|\[\]\{\}]+)(?:\:([\w\|\,]+))?$`
	_PARSE_TYPE_ENUM       = ` *\| *`
	_PARSE_TYPE_COMBINABLE = ` *\, *`
	_PARSE_NAME            = `^(\w+)(?:\=(\w+))?$`
)

// Analyse a comment to extract the keywords
func parseComment(comment string) map[string][][]string {

	result := make(map[string][][]string)

	tag_re := regexp.MustCompile(_PARSE_TAG)
	tag_block_re := regexp.MustCompile(_PARSE_TAGBLOCK)

	current_tag := TAG_DESCRIPTION
	current_fields := []string{}

	// Read the comment line by line
	lns := strings.Split(comment, "\n")
	for _, ln := range lns {
		res := tag_re.FindStringSubmatch(strings.Trim(ln, " \t"))
		// if we are reading a new keyword
		if len(res) > 1 {

			// store the previous one
			if result[current_tag] == nil {
				result[current_tag] = [][]string{}
			}
			result[current_tag] = append(result[current_tag], current_fields)

			// start a new one
			current_tag = res[1]
			current_fields = []string{}
			if len(res) == 3 && res[2] != "" {
				current_fields = regexp.MustCompile("\t+").Split(res[2], -1)
			}

		} else if res = tag_block_re.FindStringSubmatch(ln); len(res) > 1 {

			// If no new keyword, but still reading content, store the field
			current_fields = append(current_fields, res[1])

		}
	}

	// Store the last one
	if result[current_tag] == nil {
		result[current_tag] = [][]string{}
	}
	result[current_tag] = append(result[current_tag], current_fields)

	return result
}

// Parse the method
func parseMethod(str string) (method string, err error) {
	if !regexp.
		MustCompile(_PARSE_METHODS).
		MatchString(str) {
		return "", fmt.Errorf("unknown method `%s`", str)
	}
	return str, nil
}

// Parse a resource
func parseResource(str string) (resource string, eventual_method string, err error) {
	if str == "" {
		return "", "", fmt.Errorf("empty resource")
	}
	res := regexp.
		MustCompile(_PARSE_RESOURCE).
		FindStringSubmatch(str)
	if len(res) != 3 {
		return "", "", fmt.Errorf("resources should be relative and start by a /")
	}
	resource = res[2]
	eventual_method = res[1]
	return
}

// Parse the description
func parseDescription(arr []string) (description string, eventual_title string, err error) {

	var clean = func(s string, by ...string) string {
		cutset := ""
		if len(by) == 0 {
			cutset = " \t\n"
		}
		return strings.Trim(s, cutset)
	}

	// Check an eventual title
	if len(arr) >= 3 && clean(arr[0]) != "" && clean(arr[1]) == "" && clean(arr[2]) != "" {
		eventual_title = clean(arr[0])
		arr[0] = ""
	}

	// Check the line breaks
	reg := regexp.MustCompile(`( *\\)$`)
	for i, _ := range arr {
		if i > 1 && clean(arr[i-1], " \t") == "" {
			arr[i-1] = "\n\n"
		}
		arr[i] = reg.ReplaceAllString(arr[i], "\n")
	}

	description = clean(strings.Replace(strings.Join(arr, " "), "\n ", "\n", -1))

	return
}

// Parse a route param, a query param or the body
func parseParameter(arr []string, is_body bool) (p Parameter, register_types []Type, err error) {
	if len(arr) == 0 {
		return p, nil, fmt.Errorf("missing definition")
	}

	// Parse line
	line := strings.Join(arr, " ")
	arr = regexp.MustCompile(_PARSE_LINE).FindStringSubmatch(line)
	if len(arr) == 0 || arr[0] == "" {
		debug("can't parse line: %s sur %s", _PARSE_LINE, line)
		return p, nil, fmt.Errorf("wrong parameter definition")
	}

	// Define attributes
	type_name, name, description := "", "", ""
	if is_body {
		type_name, description = arr[1], arr[2]
	} else {
		type_name, name, description = arr[1], arr[2], arr[3]
	}

	// Parse type
	res := regexp.MustCompile(_PARSE_TYPE).FindStringSubmatch(type_name)
	if len(res) == 0 {
		debug("can't parse type: %s sur %s", _PARSE_TYPE, type_name)
		return p, nil, fmt.Errorf("wrong type definition")
	}
	type_name = res[1]
	type_enum := res[2]

	// Parse name
	type_default := ""
	if name != "" {
		res = regexp.MustCompile(_PARSE_NAME).FindStringSubmatch(name)
		if len(res) == 0 {
			return p, nil, fmt.Errorf("wrong name definition")
		}
		name = res[1]
		type_default = res[2]
	} else {
		name = type_name
	}

	// Create the object
	p = Parameter{
		Name:        name,
		Type:        Type(type_name),
		Description: description,
	}

	// Check possible values for type
	// i.e. string:a,b,c or int:0|1|2
	enum := []string{}
	examples := []string{}
	if type_enum != "" {
		// If enum values are not combinable
		if values := regexp.MustCompile(_PARSE_TYPE_ENUM).Split(type_enum, -1); len(values) > 1 {
			enum = values
		} else if values := regexp.MustCompile(_PARSE_TYPE_COMBINABLE).Split(type_enum, -1); len(values) > 1 {
			// If enum values are combinable
			// we display examples
			enum = values
			examples = []string{values[0], values[1], values[0] + "," + values[1]}
		}
	}

	// Format the type to RAML
	_, p.Type, register_types, err = formatType(string(p.Type))
	if err != nil {
		return
	}

	// If enum, parse to the good type
	var parseList = func(list []string, to string) (res []interface{}) {
		res = []interface{}{}
		for _, s := range list {
			var v interface{}
			var err error
			switch to {
			case "number":
				v, err = strconv.ParseFloat(s, 64)
			case "integer":
				v, err = strconv.Atoi(s)
			case "string":
				v = s
			}
			if err == nil {
				res = append(res, v)
			}
		}
		return
	}
	if len(enum) != 0 {
		p.Enum = parseList(enum, string(p.Type))
		p.Example = strings.Join(examples, ", ") //parseList(examples, string(p.Type))
	}

	// If default, parse to the good type
	if type_default != "" {
		val := parseList([]string{type_default}, string(p.Type))
		if len(val) > 0 {
			p.Default = val[0]
		}
	}

	return
}

// Parse a response
func parseResponse(arr []string) (r Response, register_types []Type, err error) {
	if len(arr) == 0 {
		return Response{}, nil, fmt.Errorf("missing definition for the response")
	}

	// Parse line
	line := strings.Join(arr, " ")
	arr = regexp.MustCompile(_PARSE_LINE).FindStringSubmatch(line)
	if len(arr) == 0 || arr[0] == "" {
		debug("can't parse line: %s sur %s", _PARSE_LINE, line)
		return r, nil, fmt.Errorf("wrong response definition")
	}

	// Define attributes
	type_name, description := arr[1], arr[2:]

	// Parse type
	res := regexp.MustCompile(_PARSE_TYPE).FindStringSubmatch(type_name)
	if len(res) == 0 {
		debug("can't parse type: %s sur %s", _PARSE_TYPE, type_name)
		return r, nil, fmt.Errorf("wrong type definition")
	}
	type_name = res[1]

	// Create the object
	r = Response{
		Type: Type(type_name),
		Description: strings.Trim(strings.Join(description, " "), ` 	`),
	}

	_, r.Type, register_types, err = formatType(string(r.Type))
	if err != nil {
		return
	}

	return
}

// Parse an example
func parseExample(att []string) (e Example, err error) {
	if len(att) == 0 {
		return Example{}, fmt.Errorf("missing definition for the example")
	}

	e = Example{}

	reg_URI := regexp.MustCompile(`^\/`)
	reg_body := regexp.MustCompile(`^\{`)
	reg_response := regexp.MustCompile(`^([0-9]+): (.+)$`)
	reg_endlineJSON := regexp.MustCompile(`^(?:\}|\]),?$`)
	currently := "description"

	tabs_nb := 1
	var tabs = func(nb int) (t string) {
		for i := 0; i < nb; i++ {
			t = t + "  "
		}
		return
	}

	for _, line := range att {
		if reg_URI.MatchString(line) {
			currently = "URI"
			e.URI = line
			continue
		} else if (currently == "description" || currently == "URI") && reg_body.MatchString(line) {
			currently = "body"
			e.Body = line
			tabs_nb = 1
			continue
		} else if reg_response.MatchString(line) {
			currently = "response"
			res := reg_response.FindStringSubmatch(line)
			code, err := strconv.Atoi(res[1])
			if err != nil || code <= 0 {
				return Example{}, fmt.Errorf("invalid http code %s in response", res[1])
			}
			e.HTTPCode = uint(code)
			e.Response = res[2]
			tabs_nb = 1
			continue
		}
		switch currently {
		case "description":
			e.Description = strings.Trim(e.Description+" "+line, " ")
		case "body":
			if reg_endlineJSON.MatchString(line) {
				tabs_nb = tabs_nb - 1
			}
			e.Body = e.Body + "\n" +
				tabs(tabs_nb) + line
			if line[len(line)-1] == '{' ||
				line[len(line)-1] == '[' {
				tabs_nb = tabs_nb + 1
			}
		case "response":
			if reg_endlineJSON.MatchString(line) {
				tabs_nb = tabs_nb - 1
			}
			e.Response = e.Response + "\n" +
				tabs(tabs_nb) + line
			if line[len(line)-1] == '{' ||
				line[len(line)-1] == '[' {
				tabs_nb = tabs_nb + 1
			}
		}
	}
	return
}

func parseTrait(arr interface{}) (t Trait, register_types string, err error) {
	// todo
	return
}
func parseSecurity(arr interface{}) (s Security, register_types string, err error) {
	// todo
	return
}
func parseAnnotation(arr interface{}) (a Annotation, register_types string, err error) {
	// todo
	return
}
