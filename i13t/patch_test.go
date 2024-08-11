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

func TestPatch(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch reqTimes {
		case 1:
			testReq := RequestFromBuffer(r.Body)
			w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON))
			w.WriteHeader(http.StatusOK)

			res := TestResponse{Num: testReq.Num}
			if err := json.NewEncoder(w).Encode(res); err != nil {
				t.Fatal(err)
			}
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
					Method: http.MethodPatch,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodPatch,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{
					"Content-Type": []string{fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON)},
				},
				Body: io.NopCloser(TestResponse{Num: 0}.Encode()),
			},
		},
	}

	ctx := context.Background()
	i := 0
	body := TestRequest{Num: 0}.Encode()
	for res, err := range r2.Patch(ctx, ts.URL, body) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestPatchWithContextCancel(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch reqTimes {
		case 1:
			testReq := RequestFromBuffer(r.Body)
			w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON))
			w.WriteHeader(http.StatusOK)

			res := TestResponse{Num: testReq.Num}
			if err := json.NewEncoder(w).Encode(res); err != nil {
				t.Fatal(err)
			}
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
					Method: http.MethodPatch,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
		{
			err: context.DeadlineExceeded,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	i := 0
	body := TestRequest{Num: 0}.Encode()
	for res, err := range r2.Patch(ctx, ts.URL, body, r2.WithInterval(3*time.Minute)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestPatchWithMaxRequestTimes(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch reqTimes {
		case 2:
			testReq := RequestFromBuffer(r.Body)
			w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON))
			w.WriteHeader(http.StatusOK)

			res := TestResponse{Num: testReq.Num}
			if err := json.NewEncoder(w).Encode(res); err != nil {
				t.Fatal(err)
			}
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
					Method: http.MethodPatch,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Request: &http.Request{
					Method: http.MethodPatch,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
	}

	ctx := context.Background()
	i := 0
	body := TestRequest{Num: 0}.Encode()
	for res, err := range r2.Patch(ctx, ts.URL, body, r2.WithMaxRequestTimes(2)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestPatchWithPeriod(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(30 * time.Millisecond)
		testReq := RequestFromBuffer(r.Body)
		w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON))
		w.WriteHeader(http.StatusOK)

		res := TestResponse{Num: testReq.Num}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			t.Fatal(err)
		}
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
	body := TestRequest{Num: 0}.Encode()
	for res, err := range r2.Patch(ctx, ts.URL, body, r2.WithPeriod(10*time.Millisecond), r2.WithMaxRequestTimes(2)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestPatchWithInterval(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch reqTimes {
		case 1:
			testReq := RequestFromBuffer(r.Body)
			w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON))
			w.WriteHeader(http.StatusOK)

			res := TestResponse{Num: testReq.Num}
			if err := json.NewEncoder(w).Encode(res); err != nil {
				t.Fatal(err)
			}
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
					Method: http.MethodPatch,
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
	body := TestRequest{Num: 0}.Encode()
	for res, err := range r2.Patch(ctx, ts.URL, body, r2.WithInterval(time.Minute), r2.WithMaxRequestTimes(3)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestPatchWithTerminationCondition(t *testing.T) {
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
					Method: http.MethodPatch,
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
					Method: http.MethodPatch,
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
	body := TestRequest{Num: 0}.Encode()
	for res, err := range r2.Patch(ctx, ts.URL, body, opts...) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestPatchWithContentType(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != r2.ContentTypeApplicationJSON {
			t.Errorf("unexpected Content-Type: %s", r.Header.Get("Content-Type"))
			return
		}
		switch reqTimes {
		case 1:
			testReq := RequestFromBuffer(r.Body)
			w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON))
			w.WriteHeader(http.StatusOK)

			res := TestResponse{Num: testReq.Num}
			if err := json.NewEncoder(w).Encode(res); err != nil {
				t.Fatal(err)
			}
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
					Method: http.MethodPatch,
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
					Method: http.MethodPatch,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{
					"Content-Type": []string{fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON)},
				},
				Body: io.NopCloser(TestResponse{Num: 0}.Encode()),
			},
		},
	}

	ctx := context.Background()
	i := 0
	body := TestRequest{Num: 0}.Encode()
	for res, err := range r2.Patch(ctx, ts.URL, body, r2.WithContentType(r2.ContentTypeApplicationJSON)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestPatchWithHeader(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") != "test" {
			t.Errorf("unexpected X-Test: %s", r.Header.Get("X-Test"))
			return
		}
		switch reqTimes {
		case 1:
			testReq := RequestFromBuffer(r.Body)
			w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON))
			w.WriteHeader(http.StatusOK)

			res := TestResponse{Num: testReq.Num}
			if err := json.NewEncoder(w).Encode(res); err != nil {
				t.Fatal(err)
			}
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
					Method: http.MethodPatch,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodPatch,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{
					"Content-Type": []string{fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON)},
				},
				Body: io.NopCloser(TestResponse{Num: 0}.Encode()),
			},
		},
	}

	ctx := context.Background()
	i := 0
	body := TestRequest{Num: 0}.Encode()
	for res, err := range r2.Patch(ctx, ts.URL, body, r2.WithHeader(http.Header{"X-Test": []string{"test"}})) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestPatchWithAspect(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testReq := RequestFromBuffer(r.Body)
		w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON))
		w.WriteHeader(http.StatusOK)

		res := TestResponse{Num: testReq.Num}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			t.Fatal(err)
		}
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	expect := []Result{
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodPatch,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{
					"Content-Type": []string{fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON)},
				},
				Body: io.NopCloser(TestResponse{Num: 1}.Encode()),
			},
		},
	}

	ctx := context.Background()
	i := 0
	body := TestRequest{Num: 0}.Encode()
	for res, err := range r2.Patch(ctx, ts.URL, body, r2.WithAspect(func(req *http.Request, do func(req *http.Request) (*http.Response, error)) (*http.Response, error) {
		testReq := RequestFromBuffer(req.Body)
		testReq.Num += 1
		req.Body = io.NopCloser(testReq.Encode())
		return do(req)
	})) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		i++
	}
}

func TestPatchWithAutoCloseResponseBody(t *testing.T) {
	reqTimes := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch reqTimes {
		case 1:
			testReq := RequestFromBuffer(r.Body)
			w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON))
			w.WriteHeader(http.StatusOK)

			res := TestResponse{Num: testReq.Num}
			if err := json.NewEncoder(w).Encode(res); err != nil {
				t.Fatal(err)
			}
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
					Method: http.MethodPatch,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{},
			},
		},
		{
			res: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					Method: http.MethodPatch,
					URL:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
				},
				Header: http.Header{
					"Content-Type": []string{fmt.Sprintf("%s; charset=utf-8", r2.ContentTypeApplicationJSON)},
				},
				Body: io.NopCloser(TestResponse{Num: 0}.Encode()),
			},
		},
	}

	ctx := context.Background()
	i := 0
	body := TestRequest{Num: 0}.Encode()
	for res, err := range r2.Patch(ctx, ts.URL, body, r2.WithAutoCloseResponseBody(false)) {
		Cmp(t, Result{res: res, err: err}, expect[i])
		if resBody := res.Body; resBody != nil {
			if err := res.Body.Close(); err != nil {
				t.Fatal(err)
			}
		}
		i++
	}
}
