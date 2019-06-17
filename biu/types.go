package biu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Method string

const (
	MethodGet    Method = "GET"
	MethodPost   Method = "POST"
	MethodPut    Method = "PUT"
	MethodPatch  Method = "PATCH"
	MethodDelete Method = "DELETE"
	//MethodHead Method = "HEAD"
	//MethodOptions Method = "OPTIONS"
	//MethodConnect Method = "CONNECT"
	//MethodTrace Method = "TRACE"
)

type Cookie = http.Cookie

type AllowRedirects struct {
	Max            uint
	Protocols      []Protocol
	Strict         bool
	Referer        bool
	TrackRedirects bool
}

type Protocol string

const (
	ProtocolHttp  Protocol = "http"
	ProtocolHttps Protocol = "https"
)

type Proxy string

const (
	ProxyHttp  Proxy = "http"
	ProxyHttps Proxy = "https"
	ProxyNo    Proxy = "no"
)

type Body map[string]interface{}

// parse request body to url-encoded form
func (body Body) form() (reader io.Reader, err error) {
	form := url.Values{}

	for k, v := range body {

		//_v, ok := v.(fmt.Stringer)
		//if !ok {
		//	return nil, errors.New("form value's type must implement fmt.Stringer")
		//}
		//form.Add(k, _v.String())
		form.Add(k, fmt.Sprintf("%v", v))
	}

	return strings.NewReader(form.Encode()), nil
}

// parse request body to json
func (body Body) json() (reader io.Reader, err error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

type UrlParam map[string]string

// construct url params
func (param UrlParam) string() string {
	params := url.Values{}

	for k, v := range param {
		params.Add(k, v)
	}

	return params.Encode()
}
