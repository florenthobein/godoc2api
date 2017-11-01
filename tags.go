// Centralize the tags that are read
// from the comments

package godoc2api

import "regexp"

// Reserved tags
const (
	TAG_HANDLER     = "handler"     // a http handler func that is executed when the route is called
	TAG_METHOD      = "method"      // HTTP method, ex: GET
	TAG_RESOURCE    = "resource"    // HTTP resource, can also contain the method, ex: GET /myroute
	TAG_DESCRIPTION = "description" // description of the route (string or []string)
	TAG_ROUTE       = "route"       // route parameter ([]string)
	TAG_ROUTES      = "routes"      // route parameters ([][]string)
	TAG_QUERY       = "query"       // query parameter ([]string)
	TAG_QUERIES     = "queries"     // query parameters ([][]string)
	TAG_BODY        = "body"        // body (string or []string)
	TAG_EXAMPLE     = "example"     // eventual example describing the use of the route ([]string)
	TAG_EXAMPLES    = "examples"    // eventual examples describing the use of the route ([][]string)
	TAG_RESPONSE    = "response"    // response type (string or []string)
)

// Tag types
const (
	_               = iota
	_TAG_TYPE_TRAIT // = iota that starts at 1
	_TAG_TYPE_SECURITY
	_TAG_TYPE_ANNOTATION
)

// Regex to match the reserved tags
const _RESERVED_TAGS = `(` +
	TAG_HANDLER + `|` +
	TAG_METHOD + `|` +
	TAG_RESOURCE + `|` +
	TAG_DESCRIPTION + `|` +
	TAG_ROUTE + `|` +
	TAG_QUERY + `|` +
	TAG_BODY + `|` +
	TAG_EXAMPLE + `|` +
	TAG_RESPONSE + `)`

// Registry of tags
var index_tag map[string]uint

// Reserve a tag
func reserveTag(s string, tag_type uint) {
	if index_tag == nil {
		index_tag = make(map[string]uint)
	}
	index_tag[s] = tag_type
}

// Verify if a tag is reserved
func isReservedTag(s string) (tag_type uint, ok bool) {
	if index_tag != nil {
		tag_type, ok = index_tag[s]
	}
	return tag_type, ok || regexp.
		MustCompile(`^`+_RESERVED_TAGS+`$`).
		MatchString(s)
}
func isReservedTrait(s string) bool {
	tag_type, ok := isReservedTag(s)
	return ok && tag_type == _TAG_TYPE_TRAIT
}
func isReservedSecurity(s string) bool {
	tag_type, ok := isReservedTag(s)
	return ok && tag_type == _TAG_TYPE_SECURITY
}
func isReservedAnnotation(s string) bool {
	tag_type, ok := isReservedTag(s)
	return ok && tag_type == _TAG_TYPE_ANNOTATION
}

// Verify if a tag type is reserved
func hasReservedTagType(tag_type uint) (ok bool) {
	if index_tag == nil {
		return
	}
	for _, typ := range index_tag {
		if typ == tag_type {
			ok = true
			break
		}
	}
	return
}

func hasReservedTrait() bool {
	return hasReservedTagType(_TAG_TYPE_TRAIT)
}
func hasReservedSecurity() bool {
	return hasReservedTagType(_TAG_TYPE_SECURITY)
}
func hasReservedAnnotation() bool {
	return hasReservedTagType(_TAG_TYPE_ANNOTATION)
}
