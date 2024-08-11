package i13t

import (
	"context"
	"fmt"
	"github.com/miyamo2/r2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestHead(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch reqTimes {
		case 1:
			w.WriteHeader(http.StatusOK)
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
					Method: http.MethodHead,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodHead,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
	}

	ctx := context.Background()
	i := 0
	for res, err := range r2.Head(ctx, ts.URL) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestHeadWithContextCancel(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(150 * time.Millisecond)
		switch reqTimes {
		case 1:
			w.WriteHeader(http.StatusOK)
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
					Method: http.MethodHead,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
		{
			err: context.DeadlineExceeded,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	i := 0
	for res, err := range r2.Head(ctx, ts.URL) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestHeadWithMaxRequestTimes(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch reqTimes {
		case 2:
			w.WriteHeader(http.StatusOK)
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
					Method: http.MethodHead,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Request: &http.Request{
					Method: http.MethodHead,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
	}

	ctx := context.Background()
	i := 0
	for res, err := range r2.Head(ctx, ts.URL, r2.WithMaxRequestTimes(2)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestHeadWithPeriod(t *testing.T) {
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
	for res, err := range r2.Head(ctx, ts.URL, r2.WithPeriod(10*time.Millisecond), r2.WithMaxRequestTimes(2)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestHeadWithInterval(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch reqTimes {
		case 1:
			w.WriteHeader(http.StatusOK)
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
					Method: http.MethodHead,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
		{
			err: context.DeadlineExceeded,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	i := 0
	for res, err := range r2.Head(ctx, ts.URL, r2.WithInterval(time.Minute), r2.WithMaxRequestTimes(3)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestHeadWithTerminationCondition(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-My-Header", fmt.Sprintf("%d", reqTimes))
		w.WriteHeader(http.StatusOK)
		reqTimes++
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	expect := []Result{
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodHead,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{
					"X-My-Header": []string{"0"},
				},
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodHead,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{
					"X-My-Header": []string{"1"},
				},
			},
		},
	}

	opts := []r2.Option{
		r2.WithContentType(r2.ContentTypeApplicationJSON),
		r2.WithTerminationCondition(func(res *http.Response) bool {
			xMyHeader := res.Header.Get("X-My-Header")
			return xMyHeader == "1"
		}),
	}

	ctx := context.Background()
	i := 0
	for res, err := range r2.Head(ctx, ts.URL, opts...) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestHeadWithContentType(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "" {
			t.Errorf("unexpected Content-Type: %s", r.Header.Get("Content-Type"))
			return
		}
		switch reqTimes {
		case 1:
			w.WriteHeader(http.StatusOK)
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
					Method: http.MethodHead,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodHead,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
	}

	ctx := context.Background()
	i := 0
	for res, err := range r2.Head(ctx, ts.URL, r2.WithContentType(r2.ContentTypeApplicationJSON)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestHeadWithHeader(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") != "test" {
			t.Errorf("unexpected X-Test: %s", r.Header.Get("X-Test"))
			return
		}
		switch reqTimes {
		case 1:
			w.WriteHeader(http.StatusOK)
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
					Method: http.MethodHead,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodHead,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
	}

	ctx := context.Background()
	i := 0
	for res, err := range r2.Head(ctx, ts.URL, r2.WithHeader(http.Header{"X-Test": []string{"test"}})) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestHeadWithAspect(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") != "test" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	expect := []Result{
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodHead,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
	}

	ctx := context.Background()
	i := 0
	for res, err := range r2.Head(ctx, ts.URL, r2.WithAspect(func(req *http.Request, do func(req *http.Request) (*http.Response, error)) (*http.Response, error) {
		req.Header.Set("X-Test", "test")
		return do(req)
	})) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}
