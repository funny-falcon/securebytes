# Secure Bytes [![GoDoc](https://godoc.org/github.com/github.com/meehow/securebytes?status.svg)](http://godoc.org/github.com/meehow/securebytes) [![Go Report Card](https://goreportcard.com/badge/github.com/meehow/securebytes)](https://goreportcard.com/report/github.com/meehow/securebytes#SecureBytes)

Secure Bytes takes any Go data type, serializes it to JSON or GOB and encrypts it with AES-192-GCM.
The goal of this library is to generate smaller cookies than [securecookie](https://github.com/gorilla/securecookie) does.
It's achieved by using [Authenticated Encryption](https://en.wikipedia.org/wiki/Authenticated_encryption) instead of HMAC.

## Installation

```
go get -u github.com/meehow/securebytes
```

## Usage

First, create Secure Bytes instance and set encryption key.
Suggested key length is at least 50 characters.

```go
var sb = securebytes.New([]byte("choo}ng-o9quoh6oodurabishoh9haeWee~neeyaRoqu6Chue1"))
```

Write a cookie:

```go
func SetCookieHandler(w http.ResponseWriter, r *http.Request) {
	secret := map[string]string{"foo": "bar"}
	b64, err := sb.EncodeToBase64(secret)
	if err != nil {
		fmt.Fprintf(w, "Encryption error: %v", err)
		return
	}
	cookie := &http.Cookie{
		Name:  "cookie-name",
		Value: b64,
		Path:  "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
```

Read a cookie:

```go
func ReadCookieHandler(w http.ResponseWriter, r *http.Request) {
	var secret map[string]string
	cookie, err := r.Cookie("cookie-name")
	if err != nil {
		fmt.Fprintf(w, "Cookie not found: %v", err)
		return
	}
	if err = sb.DecryptBase64(cookie.Value, &secret); err != nil {
		fmt.Fprintf(w, "Decryption error: %v", err)
		return
	}
	fmt.Fprintf(w, "Your secret cookie: %#v", secret)
}
```

You can also check [example](examples/cookie.go) to see full code of http server.

If you need plain `[]byte` output, you can use `Encrypt` and `Decrypt` functions instead.

You can find more information in the [documentation](https://godoc.org/github.com/meehow/securebytes#SecureBytes).

## Benchmark

SecureBytes works **2.5 times faster** than SecureCookie and generates **40% smaller** cookies.

```
goos: linux
goarch: amd64
pkg: github.com/meehow/securebytes
BenchmarkSecureBytesJSON-4    	  300000	      3810 ns/op
BenchmarkSecureBytesGOB-4     	  300000	      5476 ns/op
BenchmarkSecureCookieJSON-4   	  200000	      9959 ns/op
BenchmarkSecureCookieGOB-4    	  100000	     11807 ns/op
PASS
ok  	github.com/meehow/securebytes	6.303s
```
