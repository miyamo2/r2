# r2 - __range__ over http __request__
[![Go Reference](https://pkg.go.dev/badge/github.com/miyamo2/r2.svg)](https://pkg.go.dev/github.com/miyamo2/r2)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/miyamo2/r2)](https://img.shields.io/github/go-mod/go-version/miyamo2/r2)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/miyamo2/r2)](https://img.shields.io/github/v/release/miyamo2/r2)
[![codecov](https://codecov.io/gh/miyamo2/r2/graph/badge.svg?token=NL0BQIIAZJ)](https://codecov.io/gh/miyamo2/r2)
[![Go Report Card](https://goreportcard.com/badge/github.com/miyamo2/r2)](https://goreportcard.com/report/github.com/miyamo2/r2)
[![GitHub License](https://img.shields.io/github/license/miyamo2/r2?&color=blue)](https://img.shields.io/github/license/miyamo2/r2?&color=blue)

**r2** provides the ability to automatically iterate through Http requests with range over func.  
You don't need to be aware of how requests are repeated; you can focus only on the business logic.

## Quick Start

### Install

```sh
go get github.com/miyamo2/r2
```

### Setup `GOEXPERIMENT`

> [!IMPORTANT]
>
> If your Go project is Go 1.23 or higher, this section is not necessary.

```sh
go env -w GOEXPERIMENT=rangefunc
```

### Simple Usage

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithMaxRequestTimes(3),
	r2.WithPeriod(time.Second),
}
for res, err := range r2.Get(ctx, "https://example.com", opts...) {
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
```

<details>
    <summary>If without miyamo2/r2</summary>

```go
type requestResult struct {
    res *http.Response
    err error
}

ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()

doReq := func(ctx context.Context, url string, timeout time.Duration) requestResult {
	ctx, cancelReq := context.WithTimeout(ctx, timeout)
	defer cancelReq()

	resultCh := make(chan requestResult, 1)
	defer close(resultCh)
	go func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			resultCh <- requestResult{
				res: nil,
				err: err,
			}
		}
		res, err := http.DefaultClient.Do(req)
		resultCh <- requestResult{
			res: res,
			err: err,
		}
	}()

	select {
	case <-ctx.Done():
		return requestResult{
			nil,
			ctx.Err(),
		}
	case result := <-resultCh:
		return result
	}
}

maxRequestTimes := 3
resultCh := make(chan requestResult)
terminated := make(chan struct{})
go func() {
	for i := 0; i < maxRequestTimes; i++ {
		select {
		case <-terminated:
			close(resultCh)
			return
		}
		result := doReq(ctx, "https://example.com", time.Second)
		resultCh <- result
	}
}()

for {
	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			slog.ErrorContext(ctx, "deadline exceeded.", slog.Any("error", ctx.Err()))
		}
		break
	case result := <-resultCh:
		err := result.err
		if err != nil {
			if errors.Is(err, r2.ErrTerminatedWithClientErrorResponse) {
				slog.ErrorContext(ctx, "terminated with client error response.", slog.Any("error", err))
				break
			}
			slog.WarnContext(ctx, "something happened.", slog.Any("error", err))
			continue
		}
		res := result.res
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

close(terminated)
```
</details>

### Features

| Feature                                                                 | Description                                                                                                                                                |
|-------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [`Get`](https://github.com/miyamo2/r2?tab=readme-ov-file#get)           | Send HTTP Get requests until the [termination condition](https://github.com/miyamo2/r2?tab=readme-ov-file#termination-conditions) is satisfied.            |
| [`Head`](https://github.com/miyamo2/r2?tab=readme-ov-file#head)         | Send HTTP Head requests until the [termination condition](https://github.com/miyamo2/r2?tab=readme-ov-file#termination-conditions) is satisfied.           |
| [`Post`](https://github.com/miyamo2/r2?tab=readme-ov-file#post)         | Send HTTP Post requests until the [termination condition](https://github.com/miyamo2/r2?tab=readme-ov-file#termination-conditions) is satisfied.           |
| [`Put`](https://github.com/miyamo2/r2?tab=readme-ov-file#put)           | Send HTTP Put requests until the [termination condition](https://github.com/miyamo2/r2?tab=readme-ov-file#termination-conditions) is satisfied.            |
| [`Patch`](https://github.com/miyamo2/r2?tab=readme-ov-file#patch)       | Send HTTP Patch requests until the [termination condition](https://github.com/miyamo2/r2?tab=readme-ov-file#termination-conditions) is satisfied.          |
| [`Delete`](https://github.com/miyamo2/r2?tab=readme-ov-file#delete)     | Send HTTP Delete requests until the [termination condition](https://github.com/miyamo2/r2?tab=readme-ov-file#termination-conditions) is satisfied.         |
| [`PostForm`](https://github.com/miyamo2/r2?tab=readme-ov-file#postform) | Send HTTP Post requests with form until the [termination condition](https://github.com/miyamo2/r2?tab=readme-ov-file#termination-conditions) is satisfied. |

#### Get

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithMaxRequestTimes(3),
	r2.WithPeriod(time.Second),
}
for res, err := range r2.Get(ctx, "https://example.com", opts...) {
	// do something
}
```

#### Head

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithMaxRequestTimes(3),
	r2.WithPeriod(time.Second),
}
for res, err := range r2.Head(ctx, "https://example.com", opts...) {
	// do something
}
```

#### Post

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithMaxRequestTimes(3),
	r2.WithPeriod(time.Second),
	r2.WithContentType(r2.ContentTypeApplicationJson),
}
body := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
for res, err := range r2.Post(ctx, "https://example.com", body, opts...) {
	// do something
}
```

#### Put

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithMaxRequestTimes(3),
	r2.WithPeriod(time.Second),
	r2.WithContentType(r2.ContentTypeApplicationJson),
}
body := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
for res, err := range r2.Put(ctx, "https://example.com", body, opts...) {
	// do something
}
```

#### Patch

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithMaxRequestTimes(3),
	r2.WithPeriod(time.Second),
	r2.WithContentType(r2.ContentTypeApplicationJson),
}
body := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
for res, err := range r2.Patch(ctx, "https://example.com", body, opts...) {
	// do something
}
```

#### Delete

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithMaxRequestTimes(3),
	r2.WithPeriod(time.Second),
	r2.WithContentType(r2.ContentTypeApplicationJson),
}
body := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
for res, err := range r2.Delete(ctx, "https://example.com", body, opts...) {
	// do something
}
```

#### PostForm

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithMaxRequestTimes(3),
	r2.WithPeriod(time.Second),
	r2.WithContentType(r2.ContentTypeApplicationJson),
}
form := url.Values{"foo": []string{"bar"}}
for res, err := range r2.Post(ctx, "https://example.com", form, opts...) {
	// do something
}
```


#### Termination Conditions

- Request succeeds and termination condition is not specified.
- Condition specified in `WithTerminationCondition` is satisfied.
- Response status code is a `4xx Client Error` other than `429: Too Many Request`.
- Maximum number of requests specified in `WithMaxRequestTimes` is reached.
- Exceeds the deadline for the `context.Context` passed in the argument


### Options

**r2** provides the following request options

| Option                                                                                                  | Description                                                                                                                                                                                                               | Default              |
|---------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------|
| [`WithMaxRequestTimes`](https://github.com/miyamo2/r2?tab=readme-ov-file#withmaxrequesttimes)           | The maximum number of requests to be perform.</br>If less than or equal to 0 is specified, maximum number of requests does not apply.                                                                                     | `0`                  |
| [`WithPeriod`](https://github.com/miyamo2/r2?tab=readme-ov-file#withperiod)                             | The timeout period of the per request.</br>If less than or equal to 0 is specified, the timeout period does not apply. </br>If `http.Client.Timeout` is set, the shorter one is applied.                                  | `0`                  |
| [`WithInterval`](https://github.com/miyamo2/r2?tab=readme-ov-file#withinterval)                         | The interval between next request.</br>By default, the interval is calculated by the exponential backoff and jitter.</br>If response status code is 429(Too Many Request), the interval conforms to 'Retry-After' header. | `0`                  |
| [`WithTerminationCondition`](https://github.com/miyamo2/r2?tab=readme-ov-file#withterminationcondition) | The termination condition of the iterator that references the response.                                                                                                                                                   | `nil`                |
| [`WithHttpClient`](https://github.com/miyamo2/r2?tab=readme-ov-file#withhttpclient)                     | The client to use for requests.                                                                                                                                                                                           | `http.DefaultClient` |
| [`WithHeader`](https://github.com/miyamo2/r2?tab=readme-ov-file#withheader)                             | The custom http headers for the request.                                                                                                                                                                                  | `http.Header`(blank) |
| [`WithContentType`](https://github.com/miyamo2/r2?tab=readme-ov-file#withcontenttype)                   | The 'Content-Type' for the request.                                                                                                                                                                                       | `''`                 |
| [`WithAspect`](https://github.com/miyamo2/r2?tab=readme-ov-file#withaspect)                             | The behavior to the pre-request/post-request.                                                                                                                                                                             | -                    |

#### WithMaxRequestTimes

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithMaxRequestTimes(3),
}
for res, err := range r2.Get(ctx, "https://example.com", opts...) {
	// do something
}
```

#### WithPeriod

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithPeriod(time.Second),
}
for res, err := range r2.Get(ctx, "https://example.com", opts...) {
	// do something
}
```

#### WithInterval

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithInterval(time.Second),
}
for res, err := range r2.Get(ctx, "https://example.com", opts...) {
	// do something
}
```

#### WithTerminationCondition

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithTerminationCondition(func(res *http.Response) bool {
		myHeader := res.Header.Get("X-My-Header")
		return len(myHeader) > 0
	}),
}
for res, err := range r2.Get(ctx, "https://example.com", opts...) {
	// do something
}
```

#### WithHttpClient

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
var myHttpClient *http.Client = getMyHttpClient()
opts := []r2.Option{
	r2.WithHttpClient(myHttpClient),
}
for res, err := range r2.Get(ctx, "https://example.com", opts...) {
	// do something
}
```

#### WithHeader

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithHeader(http.Header{"X-My-Header": []string{"my-value"}}),
}
for res, err := range r2.Get(ctx, "https://example.com", opts...) {
	// do something
}
```

#### WithContentType

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithContentType("application/json"),
}
for res, err := range r2.Get(ctx, "https://example.com", opts...) {
	// do something
}
```

#### WithAspect

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
    r2.WithAspect(func(req *http.Request, do func(req *http.Request) (*http.Response, error)) (*http.Response, error) {
        res, err := do(req)
        res.StatusCode = res.StatusCode + 1
        return res, err
    }),
}
for res, err := range r2.Get(ctx, "https://example.com", opts...) {
    // do something
}
```

## For Contributors

Feel free to open a PR or an Issue.

### Tasks

We recommend that this section be run with [`xc`](https://github.com/joerdav/xc).

#### setup:deps

Install `mockgen` and `golangci-lint`.

```sh
go install go.uber.org/mock/mockgen@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

#### setup:goenv

Set `GOEXPERIMENT` to `rangefunc` if Go version is 1.22.

```sh
GOVER=$(go version)
if [[ $GOVER == *"go1.22"* ]]; then
  go env -w GOEXPERIMENT=rangefunc
fi
```

#### setup:mocks

Generate mock files.

```sh
go mod tidy
go generate ./...
```

#### lint

```sh
golangci-lint run --fix
```

#### test:unit

Run Unit Test

```sh
cd ./u6t
go test -v -coverpkg=github.com/miyamo2/r2 ./... -coverprofile=coverage.out 
```

## License

**r2** released under the [MIT License](https://github.com/miyamo2/r2/blob/main/LICENSE)
