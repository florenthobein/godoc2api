package godoc2api

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/florenthobein/godoc2api/raml"
)

// Fully describe a route
type Route struct {
	Name            string
	Method          string
	Resource        string
	Description     string
	Callback        string
	URIParameters   map[string]Parameter
	QueryParameters map[string]Parameter
	BodyParameters  map[string]Parameter
	Response        *Response
	Examples        map[string]Example
	Traits          map[string]Trait
	Securities      map[string]Security
	Annotations     map[string]Annotation

	_documentation *Documentation
}

func (r *Route) signature() string {
	return fmt.Sprint("%s %s", r.Method, r.Resource)
}

// Add and parse a tag to the Route object
func (r *Route) addTag(tag string, value interface{}) (err error) {
	kind := reflect.TypeOf(value).Kind()

	var checkString = func() (string, error) {
		v := ""
		if kind == reflect.String {
			v = value.(string)
		} else if kind == reflect.Slice {
			v = strings.Join(value.([]string), "")
		} else {
			return "", fmt.Errorf("wrong kind for the tag %s: string expected", tag)
		}
		return v, nil
	}
	var checkArray = func() ([]string, error) {
		v := []string{}
		if kind == reflect.String {
			if value.(string) == "" {
				return v, nil
			}
			v = []string{value.(string)}
		} else if kind == reflect.Slice {
			v = value.([]string)
		} else {
			return []string{}, fmt.Errorf("wrong kind for the tag %s: []string expected", tag)
		}
		return v, nil
	}
	var checkArrayArray = func() ([][]string, error) {
		v := [][]string{}
		if kind == reflect.String {
			if value.(string) == "" {
				return v, nil
			}
			v = [][]string{[]string{value.(string)}}
		} else if kind == reflect.Slice {
			v = value.([][]string)
		} else {
			return [][]string{}, fmt.Errorf("wrong kind for the tag %s: [][]string expected", tag)
		}
		return v, nil
	}

	switch tag {
	case TAG_METHOD:
		v, err := checkString()
		if err != nil || v == "" {
			return err
		}
		r.Method, err = parseMethod(v)
		return err
	case TAG_RESOURCE:
		v, err := checkString()
		if err != nil || v == "" {
			return err
		}
		method := ""
		r.Resource, method, err = parseResource(v)
		if method != "" {
			r.addTag(TAG_METHOD, method)
		}
		return err
	case TAG_DESCRIPTION:
		v, err := checkArray()
		if err != nil || len(v) == 0 {
			return err
		}
		title := ""
		r.Description, title, err = parseDescription(v)
		if title != "" {
			r.Name = title
		}
		return err
	case TAG_ROUTES:
		vs, err := checkArrayArray()
		if err != nil || len(vs) == 0 {
			return err
		}
		for _, v := range vs {
			r.addTag(TAG_ROUTE, v)
		}
		return err
	case TAG_ROUTE:
		v, err := checkArray()
		if err != nil || len(v) == 0 {
			return err
		}
		p, register_types, err := parseParameter(v, false)
		if err != nil {
			return err
		}
		if r.URIParameters == nil {
			r.URIParameters = make(map[string]Parameter)
		}
		r.URIParameters[p.Name] = p
		// Store a Type definition in the Documentation
		for _, new_t := range register_types {
			r._documentation.addType(new_t)
		}
		return err
	case TAG_QUERIES:
		vs, err := checkArrayArray()
		if err != nil || len(vs) == 0 {
			return err
		}
		for _, v := range vs {
			r.addTag(TAG_QUERY, v)
		}
		return err
	case TAG_QUERY:
		v, err := checkArray()
		if err != nil || len(v) == 0 {
			return err
		}
		p, register_types, err := parseParameter(v, false)
		if err != nil {
			return err
		}
		if r.QueryParameters == nil {
			r.QueryParameters = make(map[string]Parameter)
		}
		r.QueryParameters[p.Name] = p
		// Store a Type definition in the Documentation
		for _, new_t := range register_types {
			r._documentation.addType(new_t)
		}
		return err
	case TAG_BODY:
		v, err := checkArray()
		if err != nil || len(v) == 0 {
			return err
		}
		p, register_types, err := parseParameter(v, true)
		if err != nil {
			return err
		}
		if r.BodyParameters == nil {
			r.BodyParameters = make(map[string]Parameter)
		}
		r.BodyParameters[p.Name] = p
		// Store a Type definition in the Documentation
		for _, new_t := range register_types {
			r._documentation.addType(new_t)
		}
		break
	case TAG_RESPONSE:
		v, err := checkArray()
		if err != nil || len(v) == 0 {
			return err
		}
		resp, register_types, err := parseResponse(v)
		if err != nil {
			return err
		}
		r.Response = &resp
		// Store a Type definition in the Documentation
		for _, new_t := range register_types {
			r._documentation.addType(new_t)
		}
		break
	case TAG_EXAMPLES:
		vs, err := checkArrayArray()
		if err != nil || len(vs) == 0 {
			return err
		}
		for _, v := range vs {
			r.addTag(TAG_EXAMPLE, v)
		}
		return err
	case TAG_EXAMPLE:
		v, err := checkArray()
		if err != nil || len(v) == 0 {
			return err
		}
		e, err := parseExample(v)
		if err != nil {
			return err
		}
		if r.Examples == nil {
			r.Examples = map[string]Example{}
		}
		name := fmt.Sprintf("Example%d", len(r.Examples)+1)
		r.Examples[name] = e
		break
	default:
		if kw_type, ok := isReservedTag(tag); ok {
			switch kw_type {
			case _TAG_TYPE_ANNOTATION:
				// todo
				parseAnnotation(value)
				break
			case _TAG_TYPE_TRAIT:
				// todo
				parseTrait(value)
				break
			case _TAG_TYPE_SECURITY:
				// todo
				parseSecurity(value)
				break
			default:
				err = fmt.Errorf("unkown tag type %d for %s", kw_type, tag)
			}
		} else {
			err = fmt.Errorf("unkown tag %s", tag)
		}
	}
	return
}

func (r *Route) check() error {
	if r.Method == "" {
		return fmt.Errorf("no method found")
	}
	if r.Resource == "" {
		return fmt.Errorf("no resource found")
	}
	return nil
}

// Transform the route into a RAML structure
func (r *Route) fillToRAML(index *map[string]raml.Resource) (err error) {
	if index == nil {
		return nil
	}

	// Get the resource
	if _, ok := (*index)[r.Resource]; !ok {
		(*index)[r.Resource] = raml.Resource{URI: r.Resource}
	}
	res := (*index)[r.Resource]

	// Uri parameters
	res.URIParameters, err = r._parametersToRAML(r.URIParameters)
	if err != nil {
		return
	}

	// Method
	m, err := r._methodToRAML()
	if err != nil {
		return
	}
	switch r.Method {
	case "GET":
		res.Get = m
		break
	case "PATCH":
		res.Patch = m
		break
	case "PUT":
		res.Put = m
		break
	case "HEAD":
		res.Head = m
		break
	case "POST":
		res.Post = m
		break
	case "DELETE":
		res.Delete = m
		break
	case "OPTIONS":
		res.Options = m
		break
	default:
		return fmt.Errorf("unknown method `%s` for resource %s", r.Method, r.Resource)
	}
	if res.Methods == nil {
		res.Methods = []*raml.Method{}
	}
	res.Methods = append(res.Methods, m)

	// Traits
	// todo

	// Securities
	// todo

	// Annotations
	// todo

	(*index)[r.Resource] = res

	return nil
}

func (r *Route) _parametersToRAML(ps map[string]Parameter) (parameters map[string]raml.Type, err error) {
	if ps == nil {
		return nil, nil
	}

	parameters = make(map[string]raml.Type)
	for _, p := range ps {
		t, err := p.toRAML()
		if err != nil {
			return nil, err
		}
		parameters[p.Name] = t
	}

	return
}

func (r *Route) _methodToRAML() (*raml.Method, error) {

	queryParameters, err := r._parametersToRAML(r.QueryParameters)
	if err != nil {
		return nil, err
	}

	m := raml.Method{
		Name:            r.Name,
		Description:     r.Description,
		QueryParameters: queryParameters,
	}

	// Bodies
	if r.BodyParameters != nil {
		for _, p := range r.BodyParameters {
			m.Body = &raml.Body{
				JSON: &raml.Type{
					Type:        p.Type,
					Description: p.Description,
				},
			}

			// Examples
			if len(r.Examples) != 0 {
				(*(*m.Body).JSON).Examples = map[string]interface{}{}
				for k, e := range r.Examples {
					ex, err := e.toRAMLQuery()
					if err != nil {
						return nil, err
					}
					if ex == nil {
						continue
					}
					(*(*m.Body).JSON).Examples[k] = *ex
				}
			}
		}
	}

	// Response
	if r.Response != nil {
		resp, err := (*r.Response).toRAML()
		if err != nil {
			return nil, err
		}
		// Examples
		if len(r.Examples) != 0 {
			(*resp.Body.JSON).Examples = map[string]interface{}{}
			for k, e := range r.Examples {
				ex, err := e.toRAMLResponse()
				if err != nil {
					return nil, err
				}
				if ex == nil {
					continue
				}
				(*resp.Body.JSON).Examples[k] = *ex
			}
		}
		m.Responses = map[raml.HTTPCode]raml.Response{
			200: resp,
		}
	}

	return &m, nil
}
