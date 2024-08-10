# Advanced Usage

This section introduces more practical uses of `miyamo2/r2`.

## O11y

### With New Relic

```go
package main

import (
	"context"
	"github.com/miyamo2/r2"
	"github.com/newrelic/go-agent/v3/integrations/nrpkgerrors"
	"github.com/newrelic/go-agent/v3/newrelic"
	"net/http"
)

var nr *newrelic.Application // omit the acquisition of *newrelic.Application.

func main() {
	ctx := context.Background()
	tx := newrelic.FromContext(ctx)
	opts := []r2.Option{
		r2.WithAspect(func(req *http.Request, do func(req *http.Request) (*http.Response, error)) (*http.Response, error) {
			txn := nr.StartTransaction("request")
			seg := newrelic.StartExternalSegment(txn, req)
			defer seg.End()
			response, err := do(req)
			seg.Response = response
			return response, err
		}),
	}
	for res, err := range r2.Get(ctx, "https://example.com", opts...) {
		if err != nil {
			tx.NoticeError(nrpkgerrors.Wrap(err))
		}
		// do something with res
	}
}
```