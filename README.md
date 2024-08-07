# r2
[![Go Reference](https://pkg.go.dev/badge/github.com/miyamo2/r2.svg)](https://pkg.go.dev/github.com/miyamo2/r2)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/miyamo2/r2)](https://img.shields.io/github/go-mod/go-version/miyamo2/r2)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/miyamo2/r2)](https://img.shields.io/github/v/release/miyamo2/r2)
[![codecov](https://codecov.io/gh/miyamo2/r2/graph/badge.svg?token=NL0BQIIAZJ)](https://codecov.io/gh/miyamo2/r2)
[![Go Report Card](https://goreportcard.com/badge/github.com/miyamo2/r2)](https://goreportcard.com/report/github.com/miyamo2/r2)
[![GitHub License](https://img.shields.io/github/license/miyamo2/r2?&color=blue)](https://img.shields.io/github/license/miyamo2/r2?&color=blue)

__range__ over http __request__


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

### Usage

```go

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
go test -v -coverpkg=github.com/miyamo2/r2 -coverprofile=coverage.out ./...
```

## License

**r2** released under the [MIT License](https://github.com/miyamo2/r2/blob/main/LICENSE)
