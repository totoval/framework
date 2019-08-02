package context

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"

	"github.com/totoval/framework/helpers/zone"
)

type HttpContextor interface {
	LifeCycleContextor
	DataContextor

	RequestParamContextor
	RequestFileContextor
	RequestBindingContextor
	RequestHeaderContextor
	RequestRawContextor

	ResponseContextor
	ResponseFileContextor
	ResponseStreamContextor

	Request() *http.Request
	SetRequest(r *http.Request)
	Writer() gin.ResponseWriter
	SetWriter(w gin.ResponseWriter)
	Params() gin.Params
	Accepted() []string
	Keys() map[string]interface{}
	Errors() []*gin.Error

	Deadline() (deadline zone.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

type LifeCycleContextor interface {
	Copy() *gin.Context
	HandlerName() string
	HandlerNames() []string
	Handler() gin.HandlerFunc
	Next()
	IsAborted() bool
	Abort()
	AbortWithStatus(code int)
	AbortWithStatusJSON(code int, jsonObj interface{})
	AbortWithError(code int, err error) *gin.Error
	Error(err error) *gin.Error
}

type DataContextor interface {
	Set(key string, value interface{})
	Get(key string) (value interface{}, exists bool)
	MustGet(key string) interface{}
	GetString(key string) (s string)
	GetBool(key string) (b bool)
	GetInt(key string) (i int)
	GetInt64(key string) (i64 int64)
	GetFloat64(key string) (f64 float64)
	GetTime(key string) (t zone.Time)
	GetDuration(key string) (d zone.Duration)
	GetStringSlice(key string) (ss []string)
	GetStringMap(key string) (sm map[string]interface{})
	GetStringMapString(key string) (sms map[string]string)
	GetStringMapStringSlice(key string) (smss map[string][]string)
}

type RequestParamContextor interface {
	Param(key string) string
	Query(key string) string
	DefaultQuery(key, defaultValue string) string
	GetQuery(key string) (string, bool)
	QueryArray(key string) []string
	GetQueryArray(key string) ([]string, bool)
	QueryMap(key string) map[string]string
	GetQueryMap(key string) (map[string]string, bool)
	PostForm(key string) string
	DefaultPostForm(key, defaultValue string) string
	GetPostForm(key string) (string, bool)
	PostFormArray(key string) []string
	GetPostFormArray(key string) ([]string, bool)
	PostFormMap(key string) map[string]string
	GetPostFormMap(key string) (map[string]string, bool)
}

type RequestFileContextor interface {
	FormFile(name string) (*multipart.FileHeader, error)
	MultipartForm() (*multipart.Form, error)
	SaveUploadedFile(file *multipart.FileHeader, dst string) error
}

type RequestBindingContextor interface {
	Bind(obj interface{}) error
	BindJSON(obj interface{}) error
	BindXML(obj interface{}) error
	BindQuery(obj interface{}) error
	BindYAML(obj interface{}) error
	BindUri(obj interface{}) error
	MustBindWith(obj interface{}, b binding.Binding) error
	ShouldBind(obj interface{}) error
	ShouldBindJSON(obj interface{}) error
	ShouldBindXML(obj interface{}) error
	ShouldBindQuery(obj interface{}) error
	ShouldBindYAML(obj interface{}) error
	ShouldBindUri(obj interface{}) error
	ShouldBindWith(obj interface{}, b binding.Binding) error
	ShouldBindBodyWith(obj interface{}, bb binding.BindingBody) (err error)
}

type RequestHeaderContextor interface {
	ClientIP() string
	ContentType() string
	IsWebsocket() bool
	Status(code int)
	Header(key, value string)
	GetHeader(key string) string
}

type RequestRawContextor interface {
	GetRawData() ([]byte, error)
}

type ResponseContextor interface {
	SetAccepted(formats ...string)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	Cookie(name string) (string, error)
	Render(code int, r render.Render)
	HTML(code int, name string, obj interface{})
	IndentedJSON(code int, obj interface{})
	SecureJSON(code int, obj interface{})
	JSONP(code int, obj interface{})
	JSON(code int, obj interface{})
	AsciiJSON(code int, obj interface{})
	PureJSON(code int, obj interface{})
	XML(code int, obj interface{})
	YAML(code int, obj interface{})
	ProtoBuf(code int, obj interface{})
	String(code int, format string, values ...interface{})
	Redirect(code int, location string)
	Data(code int, contentType string, data []byte)
	DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string)
}

type ResponseFileContextor interface {
	File(filepath string)
	FileAttachment(filepath, filename string)
}

type ResponseStreamContextor interface {
	SSEvent(name string, message interface{})
	Stream(step func(w io.Writer) bool) bool
	Negotiate(code int, config gin.Negotiate)
	NegotiateFormat(offered ...string) string
}
