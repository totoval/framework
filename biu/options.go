package biu

type Options struct {
	//AllowRedirects AllowRedirects
	//HttpErrors     bool
	//DecodeContent  bool
	//Verify         bool
	Cookies []*Cookie
	//Proxy          map[Proxy]string // use `,` for delimiter
	Headers  map[string]string
	Body     *Body
	UrlParam *UrlParam
}
