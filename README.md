# r2 - __range__ over http __request__
[![Go Reference](https://pkg.go.dev/badge/github.com/miyamo2/r2.svg)](https://pkg.go.dev/github.com/miyamo2/r2)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/miyamo2/r2)](https://img.shields.io/github/go-mod/go-version/miyamo2/r2)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/miyamo2/r2)](https://img.shields.io/github/v/release/miyamo2/r2)
[![codecov](https://codecov.io/gh/miyamo2/r2/graph/badge.svg?token=NL0BQIIAZJ)](https://codecov.io/gh/miyamo2/r2)
[![Go Report Card](https://goreportcard.com/badge/github.com/miyamo2/r2)](https://goreportcard.com/report/github.com/miyamo2/r2)
[![GitHub License](https://img.shields.io/github/license/miyamo2/r2?&color=blue)](https://img.shields.io/github/license/miyamo2/r2?&color=blue)

**r2** provides the ability to automatically iterate through Http requests with range over func.

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
	// There is no need to close the response body yourself as auto closing is enabled by default.
}
```

<details>
    <summary>vs 'github.com/avast/retry-go'</summary>

```go
url := "http://example.com"
var buf []byte

ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()

opts := []retry.Option{
	retry.Attempts(3),
	retry.Context(ctx),
}

// In r2, the timeout period per request can be specified with the `WithPeriod` option.
client := http.Client{
	Timeout: time.Second,
}

err := retry.Do(
	func() error {
		res, err := client.Get(url)
		if err != nil {
			return err
		}
		if res == nil {
			return fmt.Errorf("response is nil")
		}
		if res.StatusCode >= http.StatusBadRequest && res.StatusCode < http.StatusInternalServerError {
			// In r2, client errors other than TooManyRequest are excluded from retries by default.
			return nil
		}
		if res.StatusCode >= http.StatusInternalServerError {
			// In r2, automatically retry if the server error response is returned by default.
			return fmt.Error("5xx: server error response")
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: expected %d, got %d", http.StatusOK, res.StatusCode)
		}

		// In r2, the response body is automatically closed by default.
		defer res.Body.Close()
		buf, err = io.ReadAll(res.Body)
		if err != nil {
			slog.ErrorContext(ctx, "failed to read response body.", slog.Any("error", err))
			return err
		}
		return nil
	},
	opts...,
)

if err != nil {
	// handle error
}

slog.InfoContext(ctx, "response", slog.String("response", string(buf)))
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
	r2.WithMaxRequestAttempts(3),
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
	r2.WithMaxRequestAttempts(3),
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
	r2.WithMaxRequestAttempts(3),
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
	r2.WithMaxRequestAttempts(3),
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
	r2.WithMaxRequestAttempts(3),
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
	r2.WithMaxRequestAttempts(3),
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
	r2.WithMaxRequestAttempts(3),
	r2.WithPeriod(time.Second),
	r2.WithContentType(r2.ContentTypeApplicationJson),
}
form := url.Values{"foo": []string{"bar"}}
for res, err := range r2.Post(ctx, "https://example.com", form, opts...) {
	// do something
}
```


#### Termination Conditions

- Request succeeded and no termination condition is specified by `WithTerminateIf`.
- Condition that specified in `WithTerminateIf` is satisfied.
- Response status code is a `4xx Client Error` other than `429: Too Many Request`.
- Maximum number of requests specified in `WithMaxRequestAttempts` is reached.
- Exceeds the deadline for the `context.Context` passed in the argument.
- When the for range loop is interrupted by a break.


### Options

**r2** provides the following request options

| Option                                                                                                | Description                                                                                                                                                                                                               | Default              |
|-------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------|
| [`WithMaxRequestAttempts`](https://github.com/miyamo2/r2?tab=readme-ov-file#withmaxrequesttimes)      | The maximum number of requests to be performed.</br>If less than or equal to 0 is specified, maximum number of requests does not apply.                                                                                   | `0`                  |
| [`WithPeriod`](https://github.com/miyamo2/r2?tab=readme-ov-file#withperiod)                           | The timeout period of the per request.</br>If less than or equal to 0 is specified, the timeout period does not apply. </br>If `http.Client.Timeout` is set, the shorter one is applied.                                  | `0`                  |
| [`WithInterval`](https://github.com/miyamo2/r2?tab=readme-ov-file#withinterval)                       | The interval between next request.</br>By default, the interval is calculated by the exponential backoff and jitter.</br>If response status code is 429(Too Many Request), the interval conforms to 'Retry-After' header. | `0`                  |
| [`WithTerminateIf`](https://github.com/miyamo2/r2?tab=readme-ov-file#withterminateif)                 | The termination condition of the iterator that references the response.                                                                                                                                                   | `nil`                |
| [`WithHttpClient`](https://github.com/miyamo2/r2?tab=readme-ov-file#withhttpclient)                   | The client to use for requests.                                                                                                                                                                                           | `http.DefaultClient` |
| [`WithHeader`](https://github.com/miyamo2/r2?tab=readme-ov-file#withheader)                           | The custom http headers for the request.                                                                                                                                                                                  | `http.Header`(blank) |
| [`WithContentType`](https://github.com/miyamo2/r2?tab=readme-ov-file#withcontenttype)                 | The 'Content-Type' for the request.                                                                                                                                                                                       | `''`                 |
| [`WithAspect`](https://github.com/miyamo2/r2?tab=readme-ov-file#withaspect)                           | The behavior to the pre-request/post-request.                                                                                                                                                                             | -                    |
| [`WithAutoCloseResponseBody`](https://github.com/miyamo2/r2?tab=readme-ov-file#withautocloseresponse) | Whether the response body is automatically closed.</br>By default, this setting is enabled.                                                                                                                               | `true`               |

#### WithMaxRequestAttempts

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
	r2.WithMaxRequestAttempts(3),
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

#### WithTerminateIf

```go
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
        res.StatusCode += 1
        return res, err
    }),
}
for res, err := range r2.Get(ctx, "https://example.com", opts...) {
    // do something
}
```

#### WithAutoCloseResponseBody

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()
opts := []r2.Option{
    r2.WithAutoCloseResponseBody(true),
}
for res, err := range r2.Get(ctx, "https://example.com", opts...) {
    // do something
}
```

### Advanced Usage

[Read more advanced usages](https://github.com/miyamo2/r2/blob/main/.doc/ADVANCED_USAGE.md)

## For Contributors

Feel free to open a PR or an Issue.  
However, you must promise to follow our [Code of Conduct](https://github.com/miyamo2/r2/blob/main/CODE_OF_CONDUCT.md).

### Tree

```sh
.
├ .doc/            # Documentation
├ .github/
│    └ workflows/  # GitHub Actions Workflow
├ internal/        # Internal Package; Shared with sub-packages.
└ u6t/             # Unit Test
```

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
GOVER=$(go mod graph)
if [[ $GOVER == *"go@1.22"* ]]; then
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

#### test:integration

Run Integration Test

```sh
cd ./i13t
go test -v -coverpkg=github.com/miyamo2/r2 ./... -coverprofile=coverage.out 
```

## License

**r2** released under the [MIT License](https://github.com/miyamo2/r2/blob/main/LICENSE)
