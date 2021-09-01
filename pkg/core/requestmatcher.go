package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

// RequestMatcher to match the http request
type RequestMatcher struct {
	request *http.Request
	target  *http.Request
	verbose bool

	matchOptions matchOptions
}

type matchOptions struct {
	withQuery bool
	withBody  bool
}

// NewRequestMatcher create a request matcher will match request method and request path
func NewRequestMatcher(request *http.Request) *RequestMatcher {
	return &RequestMatcher{request: request}
}

// NewVerboseRequestMatcher create a verbose request matcher will match request method and request path
func NewVerboseRequestMatcher(request *http.Request) *RequestMatcher {
	return &RequestMatcher{request: request, verbose: true}
}

// WithQuery returns a matcher with query
func (matcher *RequestMatcher) WithQuery() *RequestMatcher {
	matcher.matchOptions.withQuery = true
	return matcher
}

// WithBody returns a matcher with body
func (matcher *RequestMatcher) WithBody() *RequestMatcher {
	matcher.matchOptions.withBody = true
	return matcher
}

// Matches returns a matcher with given function
func (matcher *RequestMatcher) Matches(x interface{}) bool {
	target := x.(*http.Request)
	matcher.target = target
	request := matcher.request

	match := request.Method == target.Method &&
		request.URL.Path == target.URL.Path &&
		request.URL.Opaque == target.URL.Opaque

	if match {
		match = matchHeader(request.Header, matcher.target.Header)
	}

	if matcher.matchOptions.withQuery && match {
		match = request.URL.RawQuery == target.URL.RawQuery
	}

	if matcher.matchOptions.withBody && match {
		reqBody, _ := getStrFromReader(request)
		targetBody, _ := getStrFromReader(target)
		match = reqBody == targetBody
	}

	return match
}

func matchHeader(left, right http.Header) bool {
	if len(left) != len(right) {
		return false
	}

	for k, v := range left {
		if k == "Content-Type" { // it's hard to compare file upload cases
			continue
		}
		if tv, ok := right[k]; !ok || !reflect.DeepEqual(v, tv) {
			return false
		}
	}
	return true
}

func getStrFromReader(request *http.Request) (text string, err error) {
	reader := request.Body
	if reader == nil {
		return
	}

	if data, err := ioutil.ReadAll(reader); err == nil {
		text = string(data)

		// it could be read twice
		payload := strings.NewReader(text)
		request.Body = ioutil.NopCloser(payload)
	}
	return
}

// String returns the text of current object
func (matcher *RequestMatcher) String() string {
	return fmt.Sprintf("%v", matcher.request)
}
