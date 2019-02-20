package hclient

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"net/http"
)

// IAuth defines an interface for HTTP Authentication
type IAuth interface {
	Apply(req *http.Request)
}

// BasicAuth consists of required information for HTTP Basic Auth
type BasicAuth struct {
	Username string
	Password string
}

// NewBasicAuth creates a new BasicAuth object
func NewBasicAuth(username, password string) *BasicAuth {
	return &BasicAuth{Username: username, Password: password}
}

// Apply adds Basic Authorization header to a HTTP request
func (ba *BasicAuth) Apply(req *http.Request) {
	req.SetBasicAuth(ba.Username, ba.Password)
}

// DigestAuth consists of required information for HTTP Digest Auth
type DigestAuth struct {
	Realm     string
	Qop       string
	Method    string
	Nonce     string
	Opaque    string
	Algorithm string
	HA1       string
	Cnonce    string
	Path      string
	Nc        int16
	Username  string
	Password  string
}

// NewDigestAuth creates a new DigestAuth object
func NewDigestAuth(realm, username, password string) *DigestAuth {
	d := new(DigestAuth)
	d.Realm = realm
	d.Username = username
	d.Password = password
	d.Qop = "auth"
	d.Nonce = RandomString(32)
	d.Opaque = RandomString(32)
	d.HA1 = fmt.Sprintf("%x", md5.Sum([]byte(username+":"+realm+":"+password)))
	d.Nc = 0
	d.Cnonce = RandomString(32)

	return d
}

// Apply adds Digest Authorization header to a HTTP request
func (d *DigestAuth) Apply(req *http.Request) {
	d.Nc++
	HA2 := fmt.Sprintf("%x", md5.Sum([]byte(req.Method+":"+req.URL.RequestURI())))
	response := fmt.Sprintf("%x", md5.Sum([]byte(d.HA1+":"+d.Nonce+":"+fmt.Sprintf("%08x", d.Nc)+":"+d.Cnonce+":"+d.Qop+":"+HA2)))
	authHeader := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc=%08x, qop=%s, response="%s", opaque="%s"`,
		d.Username, d.Realm, d.Nonce, req.URL.RequestURI(), d.Cnonce, d.Nc, d.Qop, response, d.Opaque)
	req.Header.Set("Authorization", authHeader)
}

// RandomString returns random string of the given length
func RandomString(length int) string {
	const DICT = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const SIZE = 62
	randBytes := make([]byte, length)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err)
	}
	for i := range randBytes {
		randBytes[i] = DICT[randBytes[i]%SIZE]
	}
	return string(randBytes)
}
