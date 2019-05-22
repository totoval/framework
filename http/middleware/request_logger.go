package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/logs"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// before request

		// collect request data
		requestHeader := c.Request.Header
		requestData, err := c.GetRawData()
		if err != nil {
			fmt.Println(err.Error())
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestData)) // 关键点

		// collect response data
		responseWriter := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		c.Next()

		// after request

		// print request data
		log.Trace("totoval request trace", logs.Field{
			"header": requestHeader,
			"body":   string(requestData),
		})
		// print response data
		log.Trace("totoval response trace", logs.Field{
			"header": responseWriter.Header(),
			"body":   responseWriter.body.String(),
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
