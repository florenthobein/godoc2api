package godoc2api

import "github.com/florenthobein/godoc2api/raml"

type Response struct {
	Type        Type
	Description string
}

func (r *Response) toRAML() (resp raml.Response, err error) {
	resp = raml.Response{
		Body: raml.Body{
			JSON: &raml.Type{
				Type:        string(r.Type),
				Description: r.Description,
			},
		},
	}
	return
}
