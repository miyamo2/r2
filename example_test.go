package r2_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/miyamo2/r2"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

func Example() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithMaxRequestTimes(3),
		r2.WithPeriod(time.Second),
	}
	for res, err := range r2.Get(ctx, "https://127.0.0.1", opts...) {
		if err != nil {
			if errors.Is(err, r2.ErrTerminatedWithClientErrorResponse) {
				slog.ErrorContext(ctx, "terminated with client error response.", slog.Any("error", err))
				break
			}
			slog.WarnContext(ctx, "something happened.", slog.Any("error", err))
			continue
		}
		if res == nil {
			slog.WarnContext(ctx, "response is nil")
			continue
		}
		if res.StatusCode != http.StatusOK {
			io.Copy(io.Discard, res.Body)
			res.Body.Close()
		}

		buf, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			slog.ErrorContext(ctx, "failed to read response body.", slog.Any("error", err))
			continue
		}
		slog.InfoContext(ctx, "response", slog.String("response", string(buf)))
	}
}

func Example_Head() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithMaxRequestTimes(3),
		r2.WithPeriod(time.Second),
	}
	for res, err := range r2.Head(ctx, "https://127.0.0.1", opts...) {
		// do something
		_, _ = res, err
	}
}

func Example_Get() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithMaxRequestTimes(3),
		r2.WithPeriod(time.Second),
	}
	for res, err := range r2.Get(ctx, "https://127.0.0.1", opts...) {
		// do something
		_, _ = res, err
	}
}

func Example_Post() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithContentType(r2.ContentTypeApplicationJSON),
		r2.WithMaxRequestTimes(3),
		r2.WithPeriod(time.Second),
	}
	body := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
	for res, err := range r2.Post(ctx, "https://127.0.0.1", body, opts...) {
		// do something
		_, _ = res, err
	}
}

func Example_Put() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithContentType(r2.ContentTypeApplicationJSON),
		r2.WithMaxRequestTimes(3),
		r2.WithPeriod(time.Second),
	}
	body := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
	for res, err := range r2.Put(ctx, "https://127.0.0.1", body, opts...) {
		// do something
		_, _ = res, err
	}
}

func Example_Patch() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithContentType(r2.ContentTypeApplicationJSON),
		r2.WithMaxRequestTimes(3),
		r2.WithPeriod(time.Second),
	}
	body := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
	for res, err := range r2.Patch(ctx, "https://127.0.0.1", body, opts...) {
		// do something
		_, _ = res, err
	}
}

func Example_Delete() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithContentType(r2.ContentTypeApplicationJSON),
		r2.WithMaxRequestTimes(3),
		r2.WithPeriod(time.Second),
	}
	body := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
	for res, err := range r2.Delete(ctx, "https://127.0.0.1", body, opts...) {
		// do something
		_, _ = res, err
	}
}

func Example_PostForm() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithContentType(r2.ContentTypeApplicationJSON),
		r2.WithMaxRequestTimes(3),
		r2.WithPeriod(time.Second),
	}
	form := url.Values{"foo": []string{"bar"}}
	for res, err := range r2.PostForm(ctx, "https://127.0.0.1", form, opts...) {
		// do something
		_, _ = res, err
	}
}

var myHttpClient *http.Client

func Example_WithHttpClient() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithHttpClient(myHttpClient),
	}
	for res, err := range r2.Get(ctx, "https://127.0.0.1", opts...) {
		// do something
		_, _ = res, err
	}
}

func Example_WithHeader() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithHeader(http.Header{"X-My-Header": []string{"my-value"}}),
	}
	for res, err := range r2.Get(ctx, "https://127.0.0.1", opts...) {
		// do something
		_, _ = res, err
	}
}

func Example_WithInterval() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithInterval(time.Second),
	}
	for res, err := range r2.Get(ctx, "https://127.0.0.1", opts...) {
		// do something
		_, _ = res, err
	}
}

func Example_WithMaxRequestTimes() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithMaxRequestTimes(3),
	}
	for res, err := range r2.Get(ctx, "https://127.0.0.1", opts...) {
		// do something
		_, _ = res, err
	}
}

func Example_WithPeriod() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithPeriod(time.Second),
	}
	for res, err := range r2.Get(ctx, "https://127.0.0.1", opts...) {
		// do something
		_, _ = res, err
	}
}

func Example_WithTerminationCondition() {
	ctx := context.Background()
	opts := []r2.Option{
		r2.WithTerminationCondition(func(res *http.Response) bool {
			myHeader := res.Header.Get("X-My-Header")
			return len(myHeader) > 0
		}),
	}
	for res, err := range r2.Get(ctx, "https://127.0.0.1", opts...) {
		// do something
		_, _ = res, err
	}
}
