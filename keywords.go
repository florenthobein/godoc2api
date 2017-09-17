// Centralize the keywords that are read
// from the comments

package doc2raml

import "regexp"

const (
	KEYWORD_CALLBACK    = "callback"
	KEYWORD_METHOD      = "method"
	KEYWORD_RESOURCE    = "resource"
	KEYWORD_DESCRIPTION = "description"
	KEYWORD_ROUTE       = "route"
	KEYWORD_QUERY       = "query"
	KEYWORD_BODY        = "body"
	KEYWORD_EXAMPLE     = "example"
	KEYWORD_RESPONSE    = "response"

	_ = iota
	KEYWORD_TYPE_TRAIT
	KEYWORD_TYPE_SECURITY
	KEYWORD_TYPE_ANNOTATION
)

const RESERVED_KEYWORDS = `(` +
	KEYWORD_CALLBACK + `|` +
	KEYWORD_METHOD + `|` +
	KEYWORD_RESOURCE + `|` +
	KEYWORD_DESCRIPTION + `|` +
	KEYWORD_ROUTE + `|` +
	KEYWORD_QUERY + `|` +
	KEYWORD_BODY + `|` +
	KEYWORD_EXAMPLE + `|` +
	KEYWORD_RESPONSE + `)`

var index_keyword map[string]uint

func reserveKeyword(s string, kw_type uint) {
	if index_keyword == nil {
		index_keyword = make(map[string]uint)
	}
	index_keyword[s] = kw_type
}

func isReservedKeyword(s string) (kw_type uint, ok bool) {
	if index_keyword != nil {
		kw_type, ok = index_keyword[s]
	}
	return kw_type, ok || regexp.
		MustCompile(`^`+RESERVED_KEYWORDS+`$`).
		MatchString(s)
}
