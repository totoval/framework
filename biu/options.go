package biu

import "time"

type Options struct {
	//AllowRedirects AllowRedirects
	//HttpErrors     bool
	//DecodeContent  bool
	//Verify         bool
	Timeout  time.Duration
	Cookies  []*Cookie
	ProxyUrl string // http://127.0.0.1:8080
	Headers  map[string]string
	Body     *Body
	UrlParam *UrlParam
}
