package doc2raml

import "github.com/cometapp/midgar/doc2raml/raml"

type Parameter struct {
	Name        string
	Type        Type
	Description string
	Enum        []interface{}
	Example     string
}

func (p *Parameter) toRAML() (t raml.Type, err error) {
	t = raml.Type{
		Name:        p.Name,
		Type:        p.Type,
		Description: p.Description,
	}
	enum := []raml.AnyType{}
	if p.Enum != nil {
		for _, e := range p.Enum {
			enum = append(enum, e)
		}
		t.Enum = enum
	}
	// if len(p.Example) != 0 {
	// 	t.Example = p.Example
	// }
	if p.Example != "" {
		t.Example = p.Example
	}
	return
}
