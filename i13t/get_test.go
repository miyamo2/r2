package i13t

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/miyamo2/r2"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	t.Parallel()
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch reqTimes {
		case 1:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		reqTimes++
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	expect := []Result{
		{
			res: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
				Body:   io.NopCloser(bytes.NewBuffer([]byte(""))),
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{
					"Content-Type": []string{"text/plain; charset=utf-8"},
				},
				Body: io.NopCloser(bytes.NewBuffer([]byte("test"))),
			},
		},
	}

	ctx := context.Background()
	i := 0
	for res, err := range r2.Get(ctx, ts.URL) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestGetWithContextCancel(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch reqTimes {
		case 1:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		reqTimes++
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	expect := []Result{
		{
			res: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
				Body:   io.NopCloser(bytes.NewBuffer([]byte(""))),
			},
		},
		{
			err: context.DeadlineExceeded,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	i := 0
	for res, err := range r2.Get(ctx, ts.URL, r2.WithInterval(3*time.Minute)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestGetWithMaxRequestAttempts(t *testing.T) {
	t.Parallel()
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch reqTimes {
		case 2:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		reqTimes++
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	expect := []Result{
		{
			res: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
				Body:   io.NopCloser(bytes.NewBuffer([]byte(""))),
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
				Body:   io.NopCloser(bytes.NewBuffer([]byte(""))),
			},
		},
	}

	ctx := context.Background()
	i := 0
	for res, err := range r2.Get(ctx, ts.URL, r2.WithMaxRequestAttempts(2)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestGetWithPeriod(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(30 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	expect := []Result{
		{
			err: context.DeadlineExceeded,
		},
		{
			err: context.DeadlineExceeded,
		},
	}

	ctx := context.Background()
	i := 0
	for res, err := range r2.Get(ctx, ts.URL, r2.WithPeriod(10*time.Millisecond), r2.WithMaxRequestAttempts(2)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestGetWithInterval(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch reqTimes {
		case 1:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		reqTimes++
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	expect := []Result{
		{
			res: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
				Body:   io.NopCloser(bytes.NewBuffer([]byte(""))),
			},
		},
		{
			err: context.DeadlineExceeded,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	i := 0
	for res, err := range r2.Get(ctx, ts.URL, r2.WithInterval(time.Minute), r2.WithMaxRequestAttempts(3)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestGetWithTerminationCondition(t *testing.T) {
	t.Parallel()
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := TestResponse{
			Num: reqTimes,
		}

		w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON))
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(body); err != nil {
			t.Fatal(err)
		}
		reqTimes++
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	expect := []Result{
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{
					"Content-Type": []string{fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON)},
				},
				Body: io.NopCloser(TestResponse{Num: 0}.Encode()),
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{
					"Content-Type": []string{fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON)},
				},
				Body: io.NopCloser(TestResponse{Num: 1}.Encode()),
			},
		},
	}

	opts := []r2.Option{
		r2.WithContentType(r2.ContentTypeApplicationJSON),
		r2.WithTerminationCondition(func(res *http.Response) bool {
			body := TestResponse{}
			err := json.NewDecoder(res.Body).Decode(&body)
			if err != nil {
				return false
			}

			return body.Num == 1
		}),
	}

	ctx := context.Background()
	i := 0
	for res, err := range r2.Get(ctx, ts.URL, opts...) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestGetWithContentType(t *testing.T) {
	t.Parallel()
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "" {
			t.Errorf("unexpected Content-Type: %s", r.Header.Get("Content-Type"))
			return
		}
		switch reqTimes {
		case 1:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		reqTimes++
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	expect := []Result{
		{
			res: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
				Body:   io.NopCloser(bytes.NewBuffer([]byte(""))),
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{
					"Content-Type": []string{"text/plain; charset=utf-8"},
				},
				Body: io.NopCloser(bytes.NewBuffer([]byte("test"))),
			},
		},
	}

	ctx := context.Background()
	i := 0
	for res, err := range r2.Get(ctx, ts.URL, r2.WithContentType(r2.ContentTypeApplicationJSON)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestGetWithHeader(t *testing.T) {
	t.Parallel()
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") != "test" {
			t.Errorf("unexpected X-Test: %s", r.Header.Get("X-Test"))
			return
		}
		switch reqTimes {
		case 1:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		reqTimes++
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	expect := []Result{
		{
			res: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
				Body:   io.NopCloser(bytes.NewBuffer([]byte(""))),
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{
					"Content-Type": []string{"text/plain; charset=utf-8"},
				},
				Body: io.NopCloser(bytes.NewBuffer([]byte("test"))),
			},
		},
	}

	ctx := context.Background()
	i := 0
	for res, err := range r2.Get(ctx, ts.URL, r2.WithHeader(http.Header{"X-Test": []string{"test"}})) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestGetWithAspect(t *testing.T) {
	t.Parallel()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	expect := []Result{
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{
					"Content-Type": []string{"text/plain; charset=utf-8"},
				},
				Body: io.NopCloser(bytes.NewBuffer([]byte("test0"))),
			},
		},
	}

	ctx := context.Background()
	i := 0
	for res, err := range r2.Get(ctx, ts.URL, r2.WithAspect(func(req *http.Request, do func(req *http.Request) (*http.Response, error)) (*http.Response, error) {
		req.Body = io.NopCloser(bytes.NewBuffer([]byte(fmt.Sprintf("%s%d", "test", i))))
		return do(req)
	})) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}
