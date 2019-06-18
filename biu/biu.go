package biu

import (
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type biu struct {
	method  Method
	url     string
	options *Options

	constructedUrl       string
	constructedRequest   *http.Request
	constructedCookieJar *cookiejar.Jar
}

func Ready(method Method, url string, options *Options) *biu {
	if options.Headers == nil {
		options.Headers = make(map[string]string)
	}

	return &biu{
		method:  method,
		url:     url,
		options: options,
	}
}
func (b *biu) Biu() (ahh *ahh) {
	return b.request()
}

func (b *biu) request() (ahh *ahh) {

	// set cookie
	if err := b.constructCookieJar(); err != nil {
		return newErrAhh(err)
	}

	// construct url params
	b.constructUrl()

	// construct request
	if err := b.constructRequest(); err != nil {
		return newErrAhh(err)
	}

	// set header
	b.constructHeader()

	client := http.Client{
		Jar: b.constructedCookieJar,
	}

	//dd, _ := httputil.DumpRequest(b.constructedRequest, true)
	//fmt.Println(string(dd))

	response, err := client.Do(b.constructedRequest)
	if err != nil {
		return newErrAhh(err)
	}
	defer response.Body.Close()

	return newAhh(response)
}

// must be used after setDefaultContentType
func (b *biu) constructRequestBody() (reader io.Reader, err error) {
	if b.options.Body == nil {
		return nil, nil
	}

	switch b.options.Headers["Content-Type"] {
	case "application/x-www-form-urlencoded":
		return b.options.Body.form()
	case "application/json":
		return b.options.Body.json()
	default:
		return b.options.Body.form()
	}
}

func (b *biu) constructUrl() {
	b.constructedUrl = b.url
	if b.options.UrlParam != nil {
		b.constructedUrl = b.constructedUrl + "?" + b.options.UrlParam.string()
	}
}
func (b *biu) constructRequest() (err error) {
	switch b.method {
	case MethodGet:
		b.constructedRequest, err = http.NewRequest(string(b.method), b.constructedUrl, nil)
		break
	case MethodPost, MethodPut, MethodPatch, MethodDelete:
		b.setDefaultContentType()
		var reqBodyBuff io.Reader
		reqBodyBuff, err = b.constructRequestBody()
		if err != nil {
			break
		}
		b.constructedRequest, err = http.NewRequest(string(b.method), b.constructedUrl, reqBodyBuff)
		break
	default:
		err = errors.New("wrong method provided")
	}

	return err
}
func (b *biu) constructCookieJar() (err error) {
	jarPtr, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	if b.options.Cookies == nil {
		b.constructedCookieJar = jarPtr
		return nil
	}

	urlPtr, err := url.Parse(b.url)
	if err != nil {
		return err
	}

	jarPtr.SetCookies(urlPtr, b.options.Cookies)

	b.constructedCookieJar = jarPtr
	return nil
}
func (b *biu) constructHeader() {
	for k, v := range b.options.Headers {
		b.constructedRequest.Header.Set(k, v)
	}
}
func (b *biu) baseUrl() (baseUrl string, err error) {
	urlPtr, err := url.Parse(b.url)
	if err != nil {
		return "", err
	}
	return urlPtr.Scheme + "://" + urlPtr.Host, nil
}

func (b *biu) setDefaultContentType() {
	if _, ok := b.options.Headers["Content-Type"]; !ok {
		b.options.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	}
}
