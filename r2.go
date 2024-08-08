/*
Package r2 provides the ability to automatically iterate through Http requests.
*/
package r2

import (
	"bytes"
	"context"
	"errors"
	"github.com/miyamo2/r2/internal"
	"io"
	"iter"
	"log/slog"
	"math"
	"math/rand/v2"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Option specifies optional parameters to r2.
type Option = internal.Option

// ErrTerminatedWithClientErrorResponse is returned when the response status code is 4xx.
// However, in the case of 429(Too Many Request) it would not be applicable.
var ErrTerminatedWithClientErrorResponse = errors.New("4xx response")

// ContentTypes
const (
	ContentTypeApplicationJSON           = "application/json"
	ContentTypeApplicationXML            = "application/xml"
	ContentTypeApplicationFormURLEncoded = "application/x-www-form-urlencoded"
	ContentTypeMultipartFormData         = "multipart/form-data"
	ContentTypeTextPlain                 = "text/plain"
	ContentTypeTextCSV                   = "text/csv"
	ContentTypeTextHTML                  = "text/html"
	ContentTypeTextCSS                   = "text/css"
	ContentTypeTextJavaScript            = "text/javascript"
	ContentTypeApplicationJavaScript     = "application/javascript"
	ContentTypeApplicationOctetStream    = "application/octet-stream"
	ContentTypeApplicationMsgPack        = "application/x-msgpack"
	ContentTypeApplicationPDF            = "application/pdf"
	ContentTypeApplicationGzip           = "application/gzip"
	ContentTypeApplicationZip            = "application/zip"
	ContentTypeApplicationLZH            = "application/x-lzh"
	ContentTypeApplicationTar            = "application/x-tar"
	ContentTypeImageBMP                  = "image/bmp"
	ContentTypeImageGIF                  = "image/gif"
	ContentTypeImageJPEG                 = "image/jpeg"
	ContentTypeImagePNG                  = "image/png"
	ContentTypeImageSVG                  = "image/svg+xml"
	ContentTypeAudioWAV                  = "audio/wav"
	ContentTypeAudioMP3                  = "audio/mp3"
	ContentTypeVideoMPEG                 = "video/mpeg"
	ContentTypeVideoMP4                  = "video/mp4"
)

// Head sends HTTP HEAD requests until one of the following conditions is satisfied.
//   - request succeeds and termination condition is not specified
//   - condition specified in [WithTerminationCondition] is satisfied
//   - response status code is a 4xx(client error) other than 429(Too Many Request)
//   - maximum number of retries specified in [WithMaxRequestTimes] is reached
//   - exceeds the deadline for the [context.Context] passed in the argument
//
// And during which time it continues to return [http.Response] and error.
func Head(ctx context.Context, url string, options ...internal.Option) iter.Seq2[*http.Response, error] {
	return responseSeq(ctx, url, http.MethodHead, nil, options...)
}

// Get sends HTTP GET requests until one of the following conditions is satisfied.
//   - request succeeds and termination condition is not specified
//   - condition specified in [WithTerminationCondition] is satisfied
//   - response status code is a 4xx(client error) other than 429(Too Many Request)
//   - maximum number of retries specified in [WithMaxRequestTimes] is reached
//   - exceeds the deadline for the [context.Context] passed in the argument
//
// And during which time it continues to return [http.Response] and error.
func Get(ctx context.Context, url string, options ...internal.Option) iter.Seq2[*http.Response, error] {
	return responseSeq(ctx, url, http.MethodGet, nil, options...)
}

// Post sends HTTP POST requests until one of the following conditions is satisfied.
//   - request succeeds and termination condition is not specified
//   - condition specified in [WithTerminationCondition] is satisfied
//   - response status code is a 4xx(client error) other than 429(Too Many Request)
//   - maximum number of retries specified in [WithMaxRequestTimes] is reached
//   - exceeds the deadline for the [context.Context] passed in the argument
//
// And during which time it continues to return [http.Response] and error.
func Post(ctx context.Context, url string, body io.Reader, options ...internal.Option) iter.Seq2[*http.Response, error] {
	return responseSeq(ctx, url, http.MethodPost, body, options...)
}

// PostForm sends HTTP POST requests until one of the following conditions is satisfied.
//   - request succeeds and termination condition is not specified
//   - condition specified in [WithTerminationCondition] is satisfied
//   - response status code is a 4xx(client error) other than 429(Too Many Request)
//   - maximum number of retries specified in [WithMaxRequestTimes] is reached
//   - exceeds the deadline for the [context.Context] passed in the argument
//
// And during which time it continues to return [http.Response] and error.
func PostForm(ctx context.Context, url string, data url.Values, options ...internal.Option) iter.Seq2[*http.Response, error] {
	options = append(options, WithContentType(ContentTypeApplicationFormURLEncoded))
	return Post(ctx, url, strings.NewReader(data.Encode()), options...)
}

// Put sends HTTP PUT requests until one of the following conditions is satisfied.
//   - request succeeds and termination condition is not specified
//   - condition specified in [WithTerminationCondition] is satisfied
//   - response status code is a 4xx(client error) other than 429(Too Many Request)
//   - maximum number of retries specified in [WithMaxRequestTimes] is reached
//   - exceeds the deadline for the [context.Context] passed in the argument
//
// And during which time it continues to return [http.Response] and error.
func Put(ctx context.Context, url string, body io.Reader, options ...internal.Option) iter.Seq2[*http.Response, error] {
	return responseSeq(ctx, url, http.MethodPut, body, options...)
}

// Patch sends HTTP PATCH requests until one of the following conditions is satisfied.
//   - request succeeds and termination condition is not specified
//   - condition specified in [WithTerminationCondition] is satisfied
//   - response status code is a 4xx(client error) other than 429(Too Many Request)
//   - maximum number of retries specified in [WithMaxRequestTimes] is reached
//   - exceeds the deadline for the [context.Context] passed in the argument
//
// And during which time it continues to return [http.Response] and error.
func Patch(ctx context.Context, url string, body io.Reader, options ...internal.Option) iter.Seq2[*http.Response, error] {
	return responseSeq(ctx, url, http.MethodPatch, body, options...)
}

// Delete sends HTTP DELETE requests until one of the following conditions is satisfied.
//   - request succeeds and termination condition is not specified
//   - condition specified in [WithTerminationCondition] is satisfied
//   - response status code is a 4xx(client error) other than 429(Too Many Request)
//   - maximum number of retries specified in [WithMaxRequestTimes] is reached
//   - exceeds the deadline for the [context.Context] passed in the argument
//
// And during which time it continues to return [http.Response] and error.
func Delete(ctx context.Context, url string, body io.Reader, options ...internal.Option) iter.Seq2[*http.Response, error] {
	return responseSeq(ctx, url, http.MethodDelete, body, options...)
}

// WithHttpClient sets a custom HTTP client for the request.
func WithHttpClient(client internal.HttpClient) internal.Option {
	return func(p *internal.R2Prop) {
		p.SetClient(client)
	}
}

// WithContentType sets the content type for the request header.
func WithContentType(contentType string) internal.Option {
	return func(p *internal.R2Prop) {
		p.SetContentType(contentType)
	}
}

// WithHeader sets custom http headers for the request.
func WithHeader(header http.Header) internal.Option {
	return func(p *internal.R2Prop) {
		p.SetHeader(header)
	}
}

// WithMaxRequestTimes sets the maximum number of requests.
// If less than or equal to 0 is specified, maximum number of requests does not apply.
func WithMaxRequestTimes(maxRequestTimes int) internal.Option {
	return func(p *internal.R2Prop) {
		p.SetMaxRequestTimes(maxRequestTimes)
	}
}

// WithInterval sets the interval between next request.
// By default, the interval is calculated by the exponential backoff and jitter.
// If response status code is 429(Too Many Request), the interval conforms to 'Retry-After' header.
func WithInterval(interval time.Duration) internal.Option {
	return func(p *internal.R2Prop) {
		p.SetInterval(interval)
	}
}

// WithPeriod sets the timeout period for the per request.
// If less than or equal to 0 is specified, the timeout period does not apply.
func WithPeriod(period time.Duration) internal.Option {
	return func(p *internal.R2Prop) {
		p.SetPeriod(period)
	}
}

// WithTerminationCondition sets the termination condition of the iterator that references the response.
func WithTerminationCondition(terminationCondition func(res *http.Response) bool) internal.Option {
	return func(p *internal.R2Prop) {
		p.SetTerminationCondition(terminationCondition)
	}
}

// responseSeq returns a sequence of [http.Response] and error.
func responseSeq(ctx context.Context, url, method string, body io.Reader, options ...internal.Option) iter.Seq2[*http.Response, error] {
	prop := internal.NewR2Prop(options...)
	client := prop.Client()
	req, err := prop.NewRequestFunc()(method, url, body)
	if err != nil {
		return noopSeq
	}
	if header := prop.Header(); header != nil {
		req.Header = header
	}
	if contentType := prop.ContentType(); contentType != "" && !slices.Contains([]string{http.MethodGet, http.MethodHead}, method) {
		req.Header.Set("Content-Type", contentType)
	}
	maxReqTimes := prop.MaxRequestTimes()

	getBody, err := rewindBody(req)
	if err != nil {
		slog.Default().WarnContext(ctx, "[r2]: request body was impossible to rewind, so the request is performed only once.", slog.Any("error", err))
		maxReqTimes = 1
	}
	return func(yield func(*http.Response, error) bool) {
		i := 0
		for {
			res, err := requestWithTimeout(ctx, client, *req, prop.Period())
			if !yield(res, err) {
				return
			}

			wait := prop.Interval()
			if res != nil {
				switch res.StatusCode {
				case http.StatusTooManyRequests:
					retryAfter := res.Header.Get(internal.ResponseHeaderKeyRetryAfter)
					if retryAfter == "" {
						break
					}
					atoi, err := strconv.Atoi(retryAfter)
					if err != nil {
						slog.Default().ErrorContext(
							ctx,
							"[r2]: server returned an invalid 'retry-after'.",
							slog.String("url", req.URL.String()),
							slog.String("retry-after", retryAfter),
							slog.Any("error", err))
						break
					}
					wait = time.Duration(atoi)
				default:
					if res.StatusCode >= http.StatusBadRequest && res.StatusCode < http.StatusInternalServerError {
						yield(nil, errors.Join(ErrTerminatedWithClientErrorResponse, errors.New(res.Status)))
						return
					}
					if cond := prop.TerminationCondition(); cond != nil {
						if cond(res) {
							return
						}
					} else if res.StatusCode < http.StatusBadRequest {
						return
					}
				}
			}

			if wait == 0 {
				wait = backOff(i)
			}
			select {
			case <-ctx.Done():
				yield(nil, ctx.Err())
				return
			case <-time.After(wait):
				// no-op
			}
			if getBody != nil {
				if req.Body, err = getBody(); err != nil {
					slog.Default().WarnContext(ctx, "[r2]: failed to rewind request body.", slog.Any("error", err))
					return
				}
			}
			i++
			if maxReqTimes != 0 && i == maxReqTimes {
				return
			}
		}
	}
}

// requestWithTimeout sends a request with a timeout.
func requestWithTimeout(ctx context.Context, client internal.HttpClient, req http.Request, period time.Duration) (*http.Response, error) {
	cancel := func() {
		// no-op
	}
	if period > 0 {
		ctx, cancel = context.WithTimeout(ctx, period)
	}
	defer cancel()
	return client.Do(req.WithContext(ctx))
}

// backOff returns the duration of the backoff.
func backOff(i int) time.Duration {
	return rand.N[time.Duration](time.Second * time.Duration(math.Pow(2, float64(i+1))))
}

// rewindBody returns an [internal.GetBodyFunc].
func rewindBody(req *http.Request) (getBody internal.GetBodyFunc, err error) {
	if req.Body == nil {
		return func() (io.ReadCloser, error) {
			return req.Body, nil
		}, nil
	}

	if req.Body == http.NoBody {
		getBody = func() (io.ReadCloser, error) {
			return req.Body, nil
		}
		return
	}
	if req.GetBody != nil {
		getBody = req.GetBody
		return
	}
	buf := bytes.Buffer{}
	tr := io.TeeReader(req.Body, &buf)
	req.Body = io.NopCloser(&buf)

	bytes, err := io.ReadAll(tr)
	if err != nil {
		return nil, err
	}
	getBody = getBodyFromBytes(bytes)
	return
}

// getBodyFromBytes returns an [internal.GetBodyFunc] from bytes.
func getBodyFromBytes(b []byte) internal.GetBodyFunc {
	return func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(b)), nil
	}
}

func noopSeq(_ func(*http.Response, error) bool) {
	// no-op
}
