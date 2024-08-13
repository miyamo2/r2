package u6t

import (
	"bytes"
	"context"
	"errors"
	"github.com/miyamo2/r2"
	"github.com/miyamo2/r2/internal"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestPostForm(t *testing.T) {
	type param struct {
		ctx     func() context.Context
		url     string
		form    url.Values
		options []internal.Option
	}
	type want struct {
		res *http.Response
		err error
	}
	type test struct {
		param                  param
		clientParamResultPairs []clientParamResultPair
		wants                  []want
	}
	tests := map[string]test{
		"most-commonly": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(2)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseInternalServerError,
						err: ErrTest,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseOK,
					},
				},
			},
			wants: []want{
				{
					res: &ResponseInternalServerError,
					err: ErrTest,
				},
				{
					res: &ResponseOK,
				},
			},
		},
		"with-termination-condition": {
			param: param{
				ctx: context.Background,
				url: "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(2), r2.WithTerminateIf(func(res *http.Response, _ error) bool {
					if xSomething, ok := res.Header["x-something"]; ok {
						return len(xSomething) == 1 && xSomething[0] == "value"
					}
					return false
				})},
				form: url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &http.Response{
							StatusCode: http.StatusOK,
						},
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"x-something": []string{"value"}},
						},
					},
				},
			},
			wants: []want{
				{
					res: &ResponseOK,
				},
				{
					res: &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"x-something": []string{"value"}},
					},
				},
			},
		},
		"with-header": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(2), r2.WithHeader(http.Header{"x-something": []string{"value"}})},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}, "x-something": []string{"value"}},
						},
					},
					result: clientResult{
						res: &ResponseOK,
					},
				},
			},
			wants: []want{
				{
					res: &ResponseOK,
				},
			},
		},
		"with-content-type": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(2), r2.WithContentType(r2.ContentTypeApplicationJSON)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseOK,
					},
				},
			},
			wants: []want{
				{
					res: &ResponseOK,
				},
			},
		},
		"with-period": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(3), r2.WithPeriod(1 * time.Nanosecond)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{},
				},
			},
			wants: []want{
				{
					err: context.DeadlineExceeded,
				},
				{
					err: context.DeadlineExceeded,
				},
				{
					err: context.DeadlineExceeded,
				},
			},
		},
		"new-request-returns-error": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestReturningError), r2.WithMaxRequestAttempts(2)},
				form:    url.Values{"foo": []string{"bar"}},
			},
		},
		"context-cancel": {
			param: param{
				ctx: func() context.Context {
					ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
					return ctx
				},
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(2), r2.WithInterval(3 * time.Minute)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseInternalServerError,
						err: ErrTest,
					},
				},
			},
			wants: []want{
				{
					res: &ResponseInternalServerError,
					err: ErrTest,
				},
				{
					err: context.DeadlineExceeded,
				},
			},
		},
		"nil-response": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(3)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						err: ErrTest,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						err: ErrTest,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						err: ErrTest,
					},
				},
			},
			wants: []want{
				{
					err: ErrTest,
				},
				{
					err: ErrTest,
				},
				{
					err: ErrTest,
				},
			},
		},
		"too-many-request": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(2)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseTooManyRequestsWithRetryAfter,
						err: ErrTest,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseOK,
					},
				},
			},
			wants: []want{
				{
					res: &ResponseTooManyRequestsWithRetryAfter,
					err: ErrTest,
				},
				{
					res: &ResponseOK,
				},
			},
		},
		"too-many-request-without-retry-after": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(2)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseTooManyRequestsWithoutRetryAfter,
						err: ErrTest,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseOK,
					},
				},
			},
			wants: []want{
				{
					res: &ResponseTooManyRequestsWithoutRetryAfter,
					err: ErrTest,
				},
				{
					res: &ResponseOK,
				},
			},
		},
		"too-many-request-with-invalid-retry-after": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(2)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseTooManyRequestsWithInvalidRetryAfter,
						err: ErrTest,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseOK,
					},
				},
			},
			wants: []want{
				{
					res: &ResponseTooManyRequestsWithInvalidRetryAfter,
					err: ErrTest,
				},
				{
					res: &ResponseOK,
				},
			},
		},
		"client-returns-not-implemented": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(2)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseNotImplemented,
						err: ErrTest,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseOK,
					},
				},
			},
			wants: []want{
				{
					res: &ResponseNotImplemented,
					err: ErrTest,
				},
				{
					res: &ResponseOK,
				},
			},
		},
		"client-returns-399": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(2)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &Response399,
					},
				},
			},
			wants: []want{
				{
					res: &Response399,
				},
			},
		},
		"client-returns-bad-request": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(2)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseBadRequest,
					},
				},
			},
			wants: []want{
				{
					res: &ResponseBadRequest,
				},
			},
		},
		"client-returns-499": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(2)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &Response499,
						err: ErrTest,
					},
				},
			},
			wants: []want{
				{
					res: &Response499,
					err: ErrTest,
				},
			},
		},
		"with-nobody": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithNoBody), r2.WithMaxRequestAttempts(2)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   http.NoBody,
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						err: ErrTest,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   http.NoBody,
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseOK,
					},
				},
			},
			wants: []want{
				{
					err: ErrTest,
				},
				{
					res: &ResponseOK,
				},
			},
		},
		"with-valid-body-without-get-body": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithValidBodyWithoutGetBody), r2.WithMaxRequestAttempts(2)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						err: ErrTest,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseOK,
					},
				},
			},
			wants: []want{
				{
					err: ErrTest,
				},
				{
					res: &ResponseOK,
				},
			},
		},
		"with-invalid-body": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithInvalidBody), r2.WithMaxRequestAttempts(2)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   &invalidReadCloser{},
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						err: ErrTest,
					},
				},
			},
			wants: []want{
				{
					err: ErrTest,
				},
			},
		},
		"with-invalid-get-body": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithInvalidGetBody), r2.WithMaxRequestAttempts(2)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						err: ErrTest,
					},
				},
			},
			wants: []want{
				{
					err: ErrTest,
				},
			},
		},
		"with-zero-max-request-times": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithMaxRequestAttempts(0)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						err: ErrTest,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						err: ErrTest,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						err: ErrTest,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseOK,
					},
				},
			},
			wants: []want{
				{
					err: ErrTest,
				},
				{
					err: ErrTest,
				},
				{
					err: ErrTest,
				},
				{
					res: &ResponseOK,
				},
			},
		},
		"with-aspect": {
			param: param{
				ctx: context.Background,
				url: "http://example.com",
				options: []internal.Option{
					internal.WithNewRequest(stubNewRequest),
					r2.WithMaxRequestAttempts(2),
					r2.WithAspect(func(req *http.Request, do func(req *http.Request) (*http.Response, error)) (*http.Response, error) {
						req.Header.Set("x-something", "value")
						res, err := do(req)
						copiedRes := &http.Response{
							StatusCode: res.StatusCode + 1,
						}
						return copiedRes, err
					}),
				},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Header: http.Header{
								"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded},
								"X-Something":  []string{"value"},
							},
						},
					},
					result: clientResult{
						res: &ResponseInternalServerError,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Header: http.Header{
								"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded},
								"X-Something":  []string{"value"},
							},
						},
					},
					result: clientResult{
						res: &ResponseOK,
					},
				},
			},
			wants: []want{
				{
					res: &http.Response{
						StatusCode: http.StatusInternalServerError + 1,
					},
				},
				{
					res: &http.Response{
						StatusCode: http.StatusOK + 1,
					},
				},
			},
		},
		"with-auto-close-response-disable": {
			param: param{
				ctx:     context.Background,
				url:     "http://example.com",
				options: []internal.Option{internal.WithNewRequest(stubNewRequestWithForm), r2.WithAutoCloseResponseBody(false)},
				form:    url.Values{"foo": []string{"bar"}},
			},
			clientParamResultPairs: []clientParamResultPair{
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseInternalServerError,
						err: ErrTest,
					},
				},
				{
					param: clientParam{
						req: &http.Request{
							URL:    HelperMustURLParse("http://example.com"),
							Method: http.MethodPost,
							Body:   io.NopCloser(bytes.NewBuffer([]byte(`foo=bar`))),
							Header: http.Header{"Content-Type": []string{r2.ContentTypeApplicationFormURLEncoded}},
						},
					},
					result: clientResult{
						res: &ResponseOK,
					},
				},
			},
			wants: []want{
				{
					res: &ResponseInternalServerError,
					err: ErrTest,
				},
				{
					res: &ResponseOK,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHttpClient := NewMockHttpClient(ctrl)
			calls := make([]any, 0, len(tt.clientParamResultPairs))
			for _, pr := range tt.clientParamResultPairs {
				calls = append(calls, mockHttpClient.EXPECT().Do(NewRequestMatcher(pr.param.req)).DoAndReturn(
					func(req *http.Request) (*http.Response, error) {
						time.Sleep(time.Second)
						return pr.result.res, pr.result.err
					},
				))
			}
			gomock.InOrder(calls...)

			i := 0
			for res, err := range r2.PostForm(tt.param.ctx(), tt.param.url, tt.param.form, append(tt.param.options, r2.WithHttpClient(mockHttpClient))...) {
				if len(tt.wants)-1 < i {
					t.Errorf("unexpected request times. expect: %d, but: %d or more", len(tt.wants), i)
				}
				w := tt.wants[i]
				if !CmpResponse(w.res, res) {
					t.Errorf("unexpected response want: %v, got: %v:\n", w.res, res)
				}
				if !errors.Is(err, w.err) {
					t.Errorf("unexpected error want: %v, got: %v", w.err, err)
				}
				i++
			}
		})
	}
}
