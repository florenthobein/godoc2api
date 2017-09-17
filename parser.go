// The parser is responsible for analysing a comment and extracting meaningful keywords
// It also structures the keyword fields

package doc2raml

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const ALLOWED_METHODS = `(GET|HEAD|POST|PUT|DELETE|PATCH|OPTIONS)`

// Analyse a comment to extract the keywords
func parseComment(comment string) map[string][][]string {

	// desc_re := regexp.MustCompile("^(?://| ?\\*) +(.+)$")
	field_re := regexp.MustCompile("^(?://| ?\\*) @(\\w+)(?:\t+(.+))?$")
	field_block_re := regexp.MustCompile("^(?://| ?\\*)(?:[ \t]+(.*))?$")

	result := make(map[string][][]string)

	current_keyword := KEYWORD_DESCRIPTION
	current_fields := []string{}

	// Read the comment line by line
	lns := strings.Split(comment, "\n")
	for _, ln := range lns {
		res := field_re.FindStringSubmatch(strings.Trim(ln, " \t"))
		// if we are reading a new keyword
		if len(res) > 1 {

			// store the previous one
			if result[current_keyword] == nil {
				result[current_keyword] = [][]string{}
			}
			result[current_keyword] = append(result[current_keyword], current_fields)

			// start a new one
			current_keyword = res[1]
			current_fields = []string{}
			if len(res) == 3 && res[2] != "" {
				current_fields = regexp.MustCompile("\t+").Split(res[2], -1)
			}

		} else if res = field_block_re.FindStringSubmatch(ln); len(res) > 1 {

			// If no new keyword, but still reading content, store the field
			current_fields = append(current_fields, res[1])

		}
	}

	// Store the last one
	if result[current_keyword] == nil {
		result[current_keyword] = [][]string{}
	}
	result[current_keyword] = append(result[current_keyword], current_fields)

	return result
}

// Parse the method
func parseMethod(str string) (method string, err error) {
	if !regexp.
		MustCompile(ALLOWED_METHODS).
		MatchString(str) {
		return "", fmt.Errorf("unknown method %s", str)
	}
	return str, nil
}

// Parse a resource
func parseResource(str string) (resource string, eventual_method string, err error) {
	if str == "" {
		return "", "", fmt.Errorf("empty resource")
	}
	if str[0] != '/' {
		return "", "", fmt.Errorf("resources should be relative and start by a /")
	}
	r := regexp.
		MustCompile(fmt.Sprintf(`^%s (.+)$`, ALLOWED_METHODS)).
		FindStringSubmatch(str)
	if len(r) == 3 {
		return r[2], r[1], nil
	}
	return str, "", nil
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

	p = Parameter{Name: arr[0]}

	if is_body {
		if len(arr) == 1 {
			p.Type = Type(arr[0])
		} else if len(arr) == 2 {
			p.Type = Type(arr[0])
			p.Description = arr[1]
		} else if len(arr) == 3 {
			p.Type = Type(arr[1])
			p.Description = arr[2]
		}
	} else {
		if len(arr) == 1 {
			err = fmt.Errorf("missing definition for %s", p.Name)
		} else if len(arr) == 2 {
			err = fmt.Errorf("missing type definition for %s", p.Name)
			p.Description = arr[1]
		} else if len(arr) == 3 {
			p.Type = Type(arr[1])
			p.Description = arr[2]
		} else {
			err = fmt.Errorf("wrong definition for %s", p.Name)
		}
	}

	// Check possible values
	// i.e. string(a,b,c) or int(0|1|2)
	re := regexp.MustCompile(`^(.+)\((.+)\)$`)
	enum := []string{}
	examples := []string{}
	if m := re.FindStringSubmatch(string(p.Type)); len(m) > 1 {
		p.Type = Type(m[1])
		// If enum values are not combinable
		if values := regexp.MustCompile(`\|`).Split(m[2], -1); len(values) > 1 {
			enum = values
		} else if values := regexp.MustCompile(`,`).Split(m[2], -1); len(values) > 1 {
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

	return
}

// Parse a response
func parseResponse(att []string) (r Response, register_types []Type, err error) {
	if len(att) == 0 {
		return Response{}, nil, fmt.Errorf("missing definition for the response")
	}

	r = Response{Type: Type(att[0])}

	if len(att) == 2 {
		r.Description = att[1]
	} else if len(att) == 3 {
		// todo what is att[1] then?
		r.Description = att[2]
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
	reg_body := regexp.MustCompile(`^\{$`)
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
