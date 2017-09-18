package godoc2api

import (
	"fmt"

	"github.com/florenthobein/godoc2api/raml"
)

type Example struct {
	URI         string
	Body        string
	Response    string
	HTTPCode    uint
	Description string
}

func (e *Example) toRAMLQuery() (ex *raml.Example, err error) {
	if e.Body == "" {
		return
	}
	ex = &raml.Example{
		Description: e.Description,
		Value:       e.Body,
		Strict:      false,
	}
	if e.URI != "" {
		ex.Description = fmt.Sprintf("%s\n`%s`", ex.Description, e.URI)
	}
	return
}

func (e *Example) toRAMLResponse() (ex *raml.Example, err error) {
	if e.Response == "" {
		return
	}
	ex = &raml.Example{
		Description: e.Description,
		Value:       e.Response,
		Strict:      false,
	}
	if e.URI != "" {
		ex.Description = fmt.Sprintf("%s\n`%s`", ex.Description, e.URI)
	}
	return
}
