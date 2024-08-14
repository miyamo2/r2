package r2_test

import (
	"bytes"
	"context"
	"github.com/miyamo2/r2"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

func Example() {
	url := "http://example.com"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithMaxRequestAttempts(3),
		r2.WithPeriod(time.Second),
	}
	for res, err := range r2.Get(ctx, url, opts...) {
		if err != nil {
			slog.WarnContext(ctx, "something happened.", slog.Any("error", err))
			continue
		}
		if res == nil {
			slog.WarnContext(ctx, "response is nil")
			continue
		}
		if res.StatusCode != http.StatusOK {
			slog.WarnContext(ctx, "unexpected status code.", slog.Int("expect", http.StatusOK), slog.Int("got", res.StatusCode))
			continue
		}

		buf, err := io.ReadAll(res.Body)
		if err != nil {
			slog.ErrorContext(ctx, "failed to read response body.", slog.Any("error", err))
			continue
		}
		slog.InfoContext(ctx, "response", slog.String("response", string(buf)))
		// There is no need to close the response body yourself as automatic closing is enabled by default.
	}
}

func ExampleHead() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithMaxRequestAttempts(3),
		r2.WithPeriod(time.Second),
	}
	for res, err := range r2.Head(ctx, "https://example.com", opts...) {
		// do something
		_, _ = res, err
	}
}

func ExampleGet() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithMaxRequestAttempts(3),
		r2.WithPeriod(time.Second),
	}
	for res, err := range r2.Get(ctx, "https://example.com", opts...) {
		// do something
		_, _ = res, err
	}
}

func ExamplePost() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithContentType(r2.ContentTypeApplicationJSON),
		r2.WithMaxRequestAttempts(3),
		r2.WithPeriod(time.Second),
	}
	body := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
	for res, err := range r2.Post(ctx, "https://example.com", body, opts...) {
		// do something
		_, _ = res, err
	}
}

func ExamplePut() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithContentType(r2.ContentTypeApplicationJSON),
		r2.WithMaxRequestAttempts(3),
		r2.WithPeriod(time.Second),
	}
	body := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
	for res, err := range r2.Put(ctx, "https://example.com", body, opts...) {
		// do something
		_, _ = res, err
	}
}

func ExamplePatch() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithContentType(r2.ContentTypeApplicationJSON),
		r2.WithMaxRequestAttempts(3),
		r2.WithPeriod(time.Second),
	}
	body := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
	for res, err := range r2.Patch(ctx, "https://example.com", body, opts...) {
		// do something
		_, _ = res, err
	}
}

func ExampleDelete() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithContentType(r2.ContentTypeApplicationJSON),
		r2.WithMaxRequestAttempts(3),
		r2.WithPeriod(time.Second),
	}
	body := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
	for res, err := range r2.Delete(ctx, "https://example.com", body, opts...) {
		// do something
		_, _ = res, err
	}
}

func ExamplePostForm() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithContentType(r2.ContentTypeApplicationJSON),
		r2.WithMaxRequestAttempts(3),
		r2.WithPeriod(time.Second),
	}
	form := url.Values{"foo": []string{"bar"}}
	for res, err := range r2.PostForm(ctx, "https://example.com", form, opts...) {
		// do something
		_, _ = res, err
	}
}

var myHttpClient *http.Client

func ExampleWithHttpClient() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithHttpClient(myHttpClient),
	}
	for res, err := range r2.Get(ctx, "https://example.com", opts...) {
		// do something
		_, _ = res, err
	}
}

func ExampleWithHeader() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithHeader(http.Header{"X-My-Header": []string{"my-value"}}),
	}
	for res, err := range r2.Get(ctx, "https://example.com", opts...) {
		// do something
		_, _ = res, err
	}
}

func ExampleWithInterval() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithInterval(time.Second),
	}
	for res, err := range r2.Get(ctx, "https://example.com", opts...) {
		// do something
		_, _ = res, err
	}
}

func ExampleWithMaxRequestAttempts() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithMaxRequestAttempts(3),
	}
	for res, err := range r2.Get(ctx, "https://example.com", opts...) {
		// do something
		_, _ = res, err
	}
}

func ExampleWithPeriod() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithPeriod(time.Second),
	}
	for res, err := range r2.Get(ctx, "https://example.com", opts...) {
		// do something
		_, _ = res, err
	}
}

func ExampleWithTerminateIf() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithTerminateIf(func(res *http.Response, _ error) bool {
			myHeader := res.Header.Get("X-My-Header")
			return len(myHeader) > 0
		}),
	}
	for res, err := range r2.Get(ctx, "https://example.com", opts...) {
		// do something
		_, _ = res, err
	}
}

func ExampleWithAutoCloseResponseBody() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithAutoCloseResponseBody(true),
	}
	for res, err := range r2.Get(ctx, "https://example.com", opts...) {
		// do something
		_, _ = res, err
	}
}

func ExampleWithAspect() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := []r2.Option{
		r2.WithAspect(func(req *http.Request, do func(req *http.Request) (*http.Response, error)) (*http.Response, error) {
			res, err := do(req)
			res.StatusCode += 1
			return res, err
		}),
	}
	for res, err := range r2.Get(ctx, "https://example.com", opts...) {
		// do something
		_, _ = res, err
	}
}
