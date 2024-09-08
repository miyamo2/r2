package unit

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/miyamo2/r2/internal"
	"go.uber.org/mock/gomock"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

var (
	ResponseOK = http.Response{
		StatusCode: http.StatusOK,
	}
	ResponseTooManyRequestsWithRetryAfter = http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{internal.ResponseHeaderKeyRetryAfter: []string{"429"}},
	}
	ResponseTooManyRequestsWithInvalidRetryAfter = http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{internal.ResponseHeaderKeyRetryAfter: []string{"abc"}},
	}
	ResponseTooManyRequestsWithoutRetryAfter = http.Response{
		StatusCode: http.StatusTooManyRequests,
	}
	Response399 = http.Response{
		StatusCode: 399,
	}
	ResponseBadRequest = http.Response{
		StatusCode: http.StatusBadRequest,
	}
	Response499 = http.Response{
		StatusCode: 499,
	}
	ResponseInternalServerError = http.Response{
		StatusCode: http.StatusInternalServerError,
	}
	ResponseNotImplemented = http.Response{
		StatusCode: http.StatusNotImplemented,
	}
)

var ErrTest = errors.New("test error")

type clientParam struct {
	req *http.Request
}

type clientResult struct {
	res *http.Response
	err error
}

type clientParamResultPair struct {
	param  clientParam
	result clientResult
}

var cmpResponseOptions = cmp.Options{
	cmpopts.IgnoreUnexported(http.Response{}),
	cmpopts.IgnoreFields(http.Response{}, "Body"),
}

var cmpRequestOptions = cmp.Options{
	cmpopts.IgnoreUnexported(http.Request{}),
	cmpopts.IgnoreFields(http.Request{}, "Body", "GetBody"),
	cmpopts.IgnoreUnexported(bytes.Buffer{}),
}

func NewRequestMatcher(x *http.Request) gomock.Matcher {
	return &RequestMatcher{
		x: x,
	}
}

type RequestMatcher struct {
	x *http.Request
}

func (m *RequestMatcher) Matches(x interface{}) bool {
	if x, ok := x.(*http.Request); ok {
		if diff := cmp.Diff(m.x, x, cmpRequestOptions); diff != "" {
			slog.Default().Error("unexpected request", slog.Any("expect", *m.x), slog.Any("got", *x))
			return false
		}
		if bodyA, bodyB := m.x.Body, x.Body; bodyA != nil && bodyB != nil {
			switch bodyA.(type) {
			case *invalidReadCloser:
				return true
			}
			bufA, err := io.ReadAll(bodyA)
			if err != nil {
				panic(err)
			}
			bufB, err := io.ReadAll(bodyB)
			if err != nil {
				panic(err)
			}
			if diff := cmp.Diff(bufA, bufB); diff != "" {
				slog.Default().Error("unexpected request body", slog.Any("expect", string(bufA)), slog.Any("got", string(bufB)))
				return false
			}
		}
	}
	return true
}
func (m *RequestMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.x)
}

func stubNewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}

	return &http.Request{
		Method: method,
		URL:    u,
		Header: http.Header{},
		Body:   io.NopCloser(body),
		GetBody: func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewBuffer([]byte(`{"foo": "bar"}`))), nil
		},
	}, nil
}

func stubNewRequestWithForm(method, urlStr string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}

	return &http.Request{
		Method: method,
		URL:    u,
		Header: http.Header{},
		Body:   io.NopCloser(body),
		GetBody: func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))), nil
		},
	}, nil
}

var errFailedToNewRequest = errors.New("failed to new request")

func stubNewRequestReturningError(_, _ string, _ io.Reader) (*http.Request, error) {
	return nil, errFailedToNewRequest
}

func stubNewRequestWithNilBody(method, urlStr string, _ io.Reader) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}

	return &http.Request{
		Method: method,
		URL:    u,
		Header: http.Header{},
	}, nil
}

func stubNewRequestWithNoBody(method, urlStr string, _ io.Reader) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}

	return &http.Request{
		Method: method,
		URL:    u,
		Body:   http.NoBody,
		Header: http.Header{},
	}, nil
}

var errFailedToReadBody = errors.New("failed to read body")

var (
	_ io.ReadCloser = (*invalidReadCloser)(nil)
)

type invalidReadCloser struct{}

func (r *invalidReadCloser) Read(_ []byte) (n int, err error) {
	return 0, errFailedToReadBody
}

func (r *invalidReadCloser) Close() error {
	return nil
}

func stubNewRequestWithValidBodyWithoutGetBody(method, urlStr string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}

	return &http.Request{
		Method: method,
		URL:    u,
		Body:   io.NopCloser(body),
		Header: http.Header{},
	}, nil
}

func stubNewRequestWithInvalidBody(method, urlStr string, _ io.Reader) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}

	return &http.Request{
		Method: method,
		URL:    u,
		Body:   io.NopCloser(&invalidReadCloser{}),
		Header: http.Header{},
	}, nil
}

var errFailedToGetBody = errors.New("failed to get body")

func stubNewRequestWithInvalidGetBody(method, urlStr string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}

	return &http.Request{
		Method: method,
		URL:    u,
		Body:   io.NopCloser(body),
		Header: http.Header{},
		GetBody: func() (io.ReadCloser, error) {
			return nil, errFailedToGetBody
		},
	}, nil
}

func HelperMustURLParse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func CmpResponse(a, b *http.Response) bool {
	if a == nil && b == nil {
		return true
	}
	if diff := cmp.Diff(a, b, cmpResponseOptions); diff != "" {
		slog.Default().Error("unexpected response", slog.Any("expect", *a), slog.Any("got", *b))
		return false
	}
	bodyA, bodyB := a.Body, b.Body
	if bodyA == nil {
		return bodyB == nil
	}
	if bodyA == http.NoBody {
		return bodyB == http.NoBody
	}
	switch bodyA.(type) {
	case *invalidReadCloser:
		return true
	}
	bufA, err := io.ReadAll(bodyA)
	if err != nil {
		return false
	}
	bufB, err := io.ReadAll(bodyB)
	if err != nil {
		return false
	}
	if diff := cmp.Diff(bufA, bufB); diff != "" {
		slog.Default().Error("unexpected response body", slog.Any("expect", string(bufA)), slog.Any("got", string(bufB)))
		return false
	}
	return true
}
