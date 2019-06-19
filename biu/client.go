package biu

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type client struct {
	http.Client
	url     *url.URL
	method  Method
	options *Options
}

func newClient(requestUrlPtr *url.URL, method Method, opts *Options) *client {
	return &client{
		url:     requestUrlPtr,
		method:  method,
		options: opts,
		Client:  http.Client{},
	}
}

func (c *client) setCookie() error {
	jarPtr, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	if c.options == nil {
		return nil
	}

	if c.options.Cookies == nil {
		return nil
	}

	jarPtr.SetCookies(c.url, c.options.Cookies)

	c.Jar = jarPtr
	return nil
}

// support http, https, socks5 proxy, example: proxyUrl = "http://127.0.0.1:8080"
func (c *client) setProxy() {
	if c.options == nil {
		return
	}
	proxy := newProxy(c.options.ProxyUrl)
	if proxy != nil {
		c.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
	}
}

func (c *client) setTimeout() {
	if c.options == nil {
		return
	}
	c.Timeout = c.options.Timeout
}
