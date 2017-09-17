// The reader is responsible for processing and storing the comments

package doc2raml

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
)

const (
	// Reflect signature of a callback
	CALLBACK_SIGNATURE = "func(http.ResponseWriter, *http.Request)"
)

// Store all the comment strings describing a callback
// identified by by its file path and line number
var index_comment map[string]map[int]string

// Analyse a callback or a struct describing a route
// and extract its comments, and eventually extra keywords
func readComment(user_route interface{}) (c string, extra map[string]interface{}, err error) {

	var callback uintptr
	extra = map[string]interface{}{}

	v := reflect.ValueOf(user_route)
	if v.Type().String() == CALLBACK_SIGNATURE {
		// If the input is a callback
		// go directly to fetching the comment
		callback = v.Pointer()
	} else {
		// If the input is a struct, read its fields
		// to find out about the callback and if possible,
		// about extra keywords
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			tag := f.Tag.Get(TAG_NAME)
			if tag == "" {
				continue
			}
			if tag == KEYWORD_CALLBACK && v.Field(i).Type().String() == CALLBACK_SIGNATURE {
				callback = v.Field(i).Pointer()
			} else if f.Type.Kind() == reflect.String {
				extra[tag] = v.Field(i).String()
			} else if f.Type.Kind() == reflect.Bool {
				extra[tag] = v.Field(i).Bool()
			}
		}
	}

	if callback == 0 {
		return "", nil, fmt.Errorf("no callback found")
	}

	// Extract the comment of the callback
	file, line := runtime.FuncForPC(callback).FileLine(callback)
	c, err = readCommentFromFile(file, line)

	// Alert if there is no comment but continue, maybe it has
	// been defined through the extra keywords
	if c == "" {
		warn("no comments for the function: %v\n", runtime.FuncForPC(callback).Name())
	}

	return
}

// Read a file to retrieve the comment of a specific callback
func readCommentFromFile(file string, line int) (string, error) {

	// Result
	c := ""

	if index_comment == nil {
		index_comment = map[string]map[int]string{}
	}

	// Reg. exp.
	comment_reg := regexp.MustCompile(`^(//| \*|\/\*\*| \*\/)`)
	func_reg := regexp.MustCompile(`^func `)

	// If the comment is already stored for that file:line
	if index_comment[file] != nil && index_comment[file][line] != "" {
		return index_comment[file][line], nil
	}

	// Read the file
	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v\n", err)
	}

	// Read the whole file to store all the comments
	r := bufio.NewReader(f)
	content, err := readLine(r)
	i := 1
	buffer := ""
	for err == nil {
		if comment_reg.MatchString(content) {
			// store the comment in a buffer
			buffer += content + "\n"
		} else if func_reg.MatchString(content) && buffer != "" {
			// the comment is complete, store it
			if index_comment[file] == nil {
				index_comment[file] = map[int]string{}
			}
			index_comment[file][i] = buffer
			buffer = ""
		} else {
			// it's not a comment that precedes a function, clear the buffer
			buffer = ""
		}

		// this is the line we are looking for
		if i == line {
			c = index_comment[file][i]
			// no break so that all the file is already read
		}

		i++
		content, err = readLine(r)
	}

	return c, nil
}

// Read a line from the buffer
func readLine(r *bufio.Reader) (string, error) {
	var err error
	var prefix bool = true
	var line, b []byte
	for prefix && err == nil {
		line, prefix, err = r.ReadLine()
		b = append(b, line...)
	}
	return string(b), err
}
