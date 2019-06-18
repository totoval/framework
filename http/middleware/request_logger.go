package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/zone"
	"github.com/totoval/framework/logs"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// before request
		startedAt := zone.Now()

		// collect request data
		requestHeader := c.Request.Header
		requestData, err := c.GetRawData()
		if err != nil {
			fmt.Println(err.Error())
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestData)) // key point

		// collect response data
		responseWriter := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		c.Next()

		// after request

		// print request
		log.Info(c.ClientIP(), logs.Field{
			"Method":         c.Request.Method,
			"Path":           c.Request.RequestURI,
			"Proto":          c.Request.Proto,
			"Status":         responseWriter.Status(),
			"UA":             c.Request.UserAgent(),
			"Latency":        zone.Now().Sub(startedAt),
			"RequestHeader":  requestHeader,
			"RequestBody":    string(requestData),
			"ResponseHeader": responseWriter.Header(),
			"ResponseBody":   responseWriter.body.String(),
		})

		// access the status we are sending
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
