package doc2raml

import "github.com/cometapp/midgar/doc2raml/raml"

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
