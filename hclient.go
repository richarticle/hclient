// A wrapper of http.Client

package hclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// DefaultClient is the default client
var DefaultClient = New(WithTimeout(time.Second*100), WithInsecureSkipVerify())

// Client structure
type Client struct {
	http.Client
	Auth IAuth
}

// New creates a new HTTP client with options using Functional Options Pattern
func New(opts ...ClientOption) *Client {
	// Set default options
	options := ClientOptions{
		Timeout:            0,
		InsecureSkipVerify: false,
		Auth:               nil,
	}

	// Update options
	for _, o := range opts {
		o(&options)
	}

	// Create client with options
	client := &Client{
		Client: http.Client{
			Timeout: options.Timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: options.InsecureSkipVerify,
				},
			},
		},
		Auth: options.Auth,
	}

	return client
}

// ClientOption is the function type for setting options
type ClientOption func(*ClientOptions)

// ClientOptions consists of all options
type ClientOptions struct {
	Timeout            time.Duration
	InsecureSkipVerify bool
	Auth               IAuth
}

// WithTimeout set HTTP timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *ClientOptions) {
		o.Timeout = timeout
	}
}

// WithInsecureSkipVerify allows insecre HTTPS
func WithInsecureSkipVerify() ClientOption {
	return func(o *ClientOptions) {
		o.InsecureSkipVerify = true
	}
}

// WithBasicAuth sets basic auth for HTTP client
func WithBasicAuth(username, password string) ClientOption {
	return func(o *ClientOptions) {
		o.Auth = NewBasicAuth(username, password)
	}
}

// WithDigestAuth sets digest auth for HTTP client
func WithDigestAuth(realm, username, password string) ClientOption {
	return func(o *ClientOptions) {
		o.Auth = NewDigestAuth(realm, username, password)
	}
}

// SetRequest is the type of middleware functions updating requests
type SetRequest func(*http.Request)

// DoJSON for default client
func DoJSON(method string, api string, reqBody interface{}, respBody interface{}, middlewares ...SetRequest) (int, error) {
	return DefaultClient.DoJSON(method, api, reqBody, respBody, middlewares...)
}

// DoJSON sends the request and gets the response in JSON format,
// and returns status code on success
func (c *Client) DoJSON(method string, api string, reqBody interface{}, respBody interface{}, middlewares ...SetRequest) (int, error) {
	// Prepare request
	var body io.Reader
	if reqBody != nil {
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(reqBody)
		body = buf
	}
	req, err := http.NewRequest(method, api, body)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	for _, m := range middlewares {
		m(req)
	}

	// Send the request
	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Output response body if provided, discard it otherwise
	if respBody != nil {
		err = json.NewDecoder(resp.Body).Decode(respBody)
		if err != nil {
			return resp.StatusCode, fmt.Errorf("Invalid response body: %v", err)
		}
	} else {
		io.Copy(ioutil.Discard, resp.Body)
	}

	return resp.StatusCode, nil
}

// Do for default client
func Do(req *http.Request) (*http.Response, error) {
	return DefaultClient.Do(req)
}

// Do sends an HTTP request and returns an HTTP response, following
// policy (such as redirects, cookies, auth) as configured on the
// client.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Apply Auth
	if c.Auth != nil {
		c.Auth.Apply(req)
	}
	return c.Client.Do(req)
}

// Get for default client
func Get(url string) (resp *http.Response, err error) {
	return DefaultClient.Get(url)
}

// Get issues a GET to the specified URL
func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post for default client
func Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return DefaultClient.Post(url, contentType, body)
}

// Post issues a POST to the specified URL
func (c *Client) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}
