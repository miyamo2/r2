//go:generate mockgen -source=$GOFILE -destination=../u6t/mock_$GOFILE -package=u6t
package internal

import (
	"io"
	"net/http"
)

// ResponseHeaderKeyRetryAfter is the header key for Retry-After
const ResponseHeaderKeyRetryAfter = "Retry-After"

// HttpClient is an abstraction of the http.Client
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// GetBodyFunc return a new copy of http.Request#Body
type GetBodyFunc func() (io.ReadCloser, error)

// NewRequest is a function that returns a new http.Request
type NewRequest func(method, url string, body io.Reader) (*http.Request, error)
