# HClient

A simple HTTP Client library for Go.

## Usage

Import hclient into your code.

```go
import "github.com/richarticle/hclient"
```

**Simple GET**
```go
resp, err := hclient.Get("https://www.google.com")
```

**Create a new client**
```go
hc := hclient.New(WithTimeout(time.Second*100), WithInsecureSkipVerify(), WithBasicAuth("username", "password"))
resp.err := hc.Get("https://myweb.com")
```

**JSON Request**
```
req := &struct {
	Name string `json:"name"`
	Age uint32  `json:"age"`
}{}

resp := &struct {
	ErrorCode string `json:"error_code"`
}{}

status, err := hclient.DoJSON("POST", "https://myrest.com/api/v1/users", req, resp)
```

