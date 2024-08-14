## 0.2.0 - 2024-08-14

### 💥Breaking Changes

- Renamed `WithMaxRequestTimes` to `WithMaxRequestAttempts`
- Added `WithTerminateIf` as a replacement for `WithTerminationCondition`
- `ErrTerminatedWithClientErrorResponse` is no longer supported
- When iterator interruption with response status code 4xx no longer yield an error
- An error caused by canceling a context is now not yielded

### 📚Documentation

- Explicitly stated that requests are terminated when the for loop is interrupted by break

## 0.1.1 - 2024-08-11

### 📚Documentation

- Fix method signatures in Examples

## 0.1.0 - 2024-08-11

### 🎉Initial Release

#### Features

##### `Get` 

Send HTTP Get requests until the termination condition is satisfied

##### `Head`  

Send HTTP Head requests until the termination condition is satisfied

##### `Post`  

Send HTTP Post requests until the termination condition is satisfied

##### `Put`  

Send HTTP Put requests until the termination condition is satisfied

##### `Patch`  

Send HTTP Patch requests until the termination condition is satisfied

##### `Delete`  

Send HTTP Delete requests until the termination condition is satisfied


##### `PostForm`  

Send HTTP Post requests with form until the termination condition is satisfied

#### Options

##### `WithMaxRequestTimes`

The maximum number of requests to be performed

##### `WithPeriod`

The timeout period of the per request

##### `WithInterval`

The interval between next request

##### `WithTerminationCondition`

The termination condition of the iterator that references the response

##### `WithHttpClient`

The client to use for requests

##### `WithHeader`

The custom http headers for the request

##### `WithContentType`

The 'Content-Type' for the request

##### `WithAspect`

The behavior to the pre-request/post-request

##### `WithAutoCloseResponseBody`

Whether the response body is automatically closed
