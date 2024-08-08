package internal

import (
	"net/http"
	"time"
)

// Option is a functional option for r2.
type Option func(*R2Prop)

// R2Prop is the properties of r2.
type R2Prop struct {
	client               HttpClient
	contentType          string
	header               http.Header
	maxRequestTimes      int
	interval             time.Duration
	period               time.Duration
	terminationCondition func(res *http.Response) bool
	newRequest           NewRequest
}

// SetClient sets the client.
func (p *R2Prop) SetClient(client HttpClient) {
	p.client = client
}

// SetContentType sets the content type.
func (p *R2Prop) SetContentType(contentType string) {
	p.contentType = contentType
}

// SetHeader sets the header.
func (p *R2Prop) SetHeader(header http.Header) {
	p.header = header
}

func (p *R2Prop) SetMaxRequestTimes(maxRequestTimes int) {
	p.maxRequestTimes = maxRequestTimes
}

// SetInterval sets the interval.
func (p *R2Prop) SetInterval(interval time.Duration) {
	p.interval = interval
}

// SetPeriod sets the period.
func (p *R2Prop) SetPeriod(period time.Duration) {
	p.period = period
}

// SetTerminationCondition sets the termination condition.
func (p *R2Prop) SetTerminationCondition(terminationCondition func(res *http.Response) bool) {
	p.terminationCondition = terminationCondition
}

// SetNewRequestFunc sets the new request function.
func (p *R2Prop) SetNewRequestFunc(newRequest NewRequest) {
	p.newRequest = newRequest
}

// Client returns the client. If the client is nil, it returns http.DefaultClient.
func (p *R2Prop) Client() HttpClient {
	if p.client == nil {
		return http.DefaultClient
	}
	return p.client
}

// Header returns the header.
func (p *R2Prop) Header() http.Header {
	return p.header
}

// MaxRequestTimes returns the max request times. If the max request times is less than or equal to 0, it returns 0.
func (p *R2Prop) MaxRequestTimes() int {
	if p.maxRequestTimes <= 0 {
		return 0
	}
	return p.maxRequestTimes
}

// Interval returns the interval. If the interval is less than 0, it returns 0.
func (p *R2Prop) Interval() time.Duration {
	if p.interval < 0 {
		return 0
	}
	return p.interval
}

// Period returns the period. If the period is less than 0, it returns 0.
func (p *R2Prop) Period() time.Duration {
	if p.period < 0 {
		return 0
	}
	return p.period
}

// NewRequestFunc returns the new request function.
func (p *R2Prop) NewRequestFunc() NewRequest {
	return p.newRequest
}

// ContentType returns the content type.
func (p *R2Prop) ContentType() string {
	return p.contentType
}

// TerminationCondition returns the termination condition.
func (p *R2Prop) TerminationCondition() func(res *http.Response) bool {
	return p.terminationCondition
}

// NewR2Prop returns a new R2Prop.
func NewR2Prop(opts ...Option) R2Prop {
	p := R2Prop{}
	for _, o := range opts {
		o(&p)
	}
	return p
}

// WithNewRequest just for testing purposes.
func WithNewRequest(newRequestFunc NewRequest) Option {
	return func(p *R2Prop) {
		p.newRequest = newRequestFunc
	}
}
