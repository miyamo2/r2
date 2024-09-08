package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"io"
	"net/http"
	"slices"
	"testing"
)

type Result struct {
	res *http.Response
	err error
}

func Cmp(t *testing.T, got, want Result) {
	t.Helper()
	if !errors.Is(got.err, want.err) {
		t.Errorf("error got: %v, want: %v", got.err, want.err)
	}
	if want.res != nil && got.res != nil {
		CmpResponse(t, got.res, want.res)
	}
}

var HeaderIgnoreEntries = []string{"Date", "Content-Length"}

var HeaderCmpOpts = []cmp.Option{
	cmpopts.IgnoreMapEntries(func(k string, _ any) bool {
		return slices.Contains(HeaderIgnoreEntries, k)
	}),
}

func CmpResponse(t *testing.T, got, want *http.Response) {
	t.Helper()
	if got.StatusCode != want.StatusCode {
		t.Errorf("StatusCode got: %d, want: %d", got.StatusCode, want.StatusCode)
	}
	if got.Request.URL.String() != want.Request.URL.String() {
		t.Errorf("Request.URL got: %s, want: %s", got.Request.URL.String(), want.Request.URL.String())
	}

	if got.Request.Method != want.Request.Method {
		t.Errorf("Request.Method got: %s, want: %s", got.Request.Method, want.Request.Method)
	}

	if diff := cmp.Diff(want.Header, got.Header, HeaderCmpOpts...); diff != "" {
		t.Errorf("Header (-want +got): %s", diff)
	}

	if diff := cmp.Diff(want.Cookies(), got.Cookies()); diff != "" {
		t.Errorf("Cookies (-want +got): %s", diff)
	}

	if wantBody, gotBody := want.Body, got.Body; wantBody != nil && gotBody != nil {
		wantBuf, err := io.ReadAll(wantBody)
		if err != nil {
			panic(err)
		}
		gotBuf, err := io.ReadAll(gotBody)
		if err != nil {
			panic(err)
		}
		if diff := cmp.Diff(wantBuf, gotBuf); diff != "" {
			t.Errorf("Body (-want +got): %s", diff)
		}
	}
}

type TestRequest struct {
	Num int `json:"num"`
}

func (r TestRequest) Encode() *bytes.Buffer {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&r); err != nil {
		return nil
	}
	return bytes.NewBuffer(buf.Bytes())
}

func RequestFromBuffer(buf io.ReadCloser) TestRequest {
	var r TestRequest
	if err := json.NewDecoder(buf).Decode(&r); err != nil {
		panic(err)
	}
	return r
}

type TestResponse struct {
	Num int `json:"num"`
}

func (r TestResponse) Encode() *bytes.Buffer {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&r); err != nil {
		return nil
	}
	return bytes.NewBuffer(buf.Bytes())
}
