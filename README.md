# Intro

Provide a consistent healthz endpoint in your services.

# Example

To generage the example response belong we have some code like the following:

```go
func main() {
	var api pkg.SomeAPI

	health := healthz.New(healthz.Meta{"service": "example"})

	m := healthz.Meta{
		"base": "http://localhost",
		"user": "root",
	}

	health.Add(healthz.NewComponent("remote_api", true, m, func() *healthz.Error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return healthz.WrapError(
			api.Ping(ctx),
			"remote_api api is not responsive",
			healthz.Meta{},
		)
	}))

	// Setup server on health.ServeHTTP
}
```

# Response Structure

For status codes 200 and 500 you'll get a response like the following.

```json
{
	"status": "unhealthy",
	"components": {
		"remote_api": {
			"healthy": false,
			"metadata": {},
			"errors": [{
				"description": "remote_api is not responsive (err: expected 200 response got 401)",
				"metadata": {
					"base": "http://localhost",
					"user": "root"
				}
			}]
		}
	},
	"metadata": {
		"service": "example"
	}
}
```

# Status Codes

## 200

It is OK to send traffic, the service may be degraded but can still provided business value.

## 500

The service is current unhealthy, a dependency might not be available, check the `components` key in the JSON response.

## 503

The service is currenting either starting or closing, check the `status` key in the JSON response.


# CLI Tool

## Installation from source.

```
$ go get -u github.com/bweston92/healthz/cmd/healthz
$ go install github.com/bweston92/healthz/cmd/healthz
$ healthz -h
```

# License

Copyright 2017 Bradley Weston

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
