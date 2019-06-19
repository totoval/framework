package biu

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type biu struct {
	method  Method
	url     *url.URL
	options *Options

	err error

	client             *client
	constructedUrl     *url.URL
	constructedRequest *http.Request
}

func Ready(method Method, requestUrl string, options *Options) *biu {
	if options != nil {
		if options.Headers == nil {
			options.Headers = make(map[string]string)
		}
	}

	requestUrlPtr, err := url.Parse(requestUrl)
	if err != nil {
		return &biu{
			method:  method,
			url:     nil,
			options: nil,
			err:     err,
		}
	}

	return &biu{
		method:  method,
		url:     requestUrlPtr,
		options: options,
	}
}
func (b *biu) Biu() (ahh *ahh) {
	if err := b.construct(); err != nil {
		return newErrAhh(err)
	}
	return b.request()
}

func (b *biu) construct() error {
	// ready is error
	if b.err != nil {
		return b.err
	}

	// construct url params
	if err := b.constructUrl(); err != nil {
		return err
	}

	// construct request
	if err := b.constructRequest(); err != nil {
		return err
	}

	// set header
	b.constructHeader()

	b.client = newClient(b.constructedUrl, b.method, b.options)

	// set cookie
	if err := b.client.setCookie(); err != nil {
		return err
	}

	// set proxy
	b.client.setProxy()

	// set timeout
	b.client.setTimeout()

	//@todo set checkRedirect func

	return nil
}

func (b *biu) request() (ahh *ahh) {

	//dd, _ := httputil.DumpRequest(b.constructedRequest, true)
	//fmt.Println(string(dd))

	defer func() {
		if err := recover(); err != nil {
			if _err, ok := err.(error); ok {
				ahh = newErrAhh(_err)
				return
			}
			ahh = newErrAhh(errors.New(fmt.Sprint(err)))
			return
		}
	}()
	response, err := b.client.Do(b.constructedRequest)
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

func (b *biu) constructUrl() (err error) {
	constructedUrl := b.url.String()
	if b.options == nil {
		b.constructedUrl = b.url
		return nil
	}
	if b.options.UrlParam != nil {
		constructedUrl = constructedUrl + "?" + b.options.UrlParam.string()
	}

	b.constructedUrl, err = url.Parse(constructedUrl)
	return err
}
func (b *biu) constructRequest() (err error) {
	switch b.method {
	case MethodGet:
		b.constructedRequest, err = http.NewRequest(string(b.method), b.constructedUrl.String(), nil)
		break
	case MethodPost, MethodPut, MethodPatch, MethodDelete:
		b.setDefaultContentType()
		var reqBodyBuff io.Reader
		reqBodyBuff, err = b.constructRequestBody()
		if err != nil {
			break
		}
		b.constructedRequest, err = http.NewRequest(string(b.method), b.constructedUrl.String(), reqBodyBuff)
		break
	default:
		err = errors.New("wrong method provided")
	}

	return err
}

func (b *biu) constructHeader() {
	if b.options == nil {
		return
	}
	for k, v := range b.options.Headers {
		b.constructedRequest.Header.Set(k, v)
	}
}
func (b *biu) baseUrl() (baseUrl string) {
	return b.url.Scheme + "://" + b.url.Host
}

func (b *biu) setDefaultContentType() {
	if b.options == nil {
		return
	}
	if _, ok := b.options.Headers["Content-Type"]; !ok {
		b.options.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	}
}
