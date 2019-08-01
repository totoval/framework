package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/config"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/helpers/zone"
	"github.com/totoval/framework/monitor/app/logics/dashboard"
	"github.com/totoval/framework/request"
)

func RequestLogger() request.HandlerFunc {
	return func(c request.Context) {

		// before request
		dashboard.Flow.Add()       // request statistics
		defer dashboard.Flow.Dec() // request statistics

		if config.GetBool("app.debug") {
			startedAt := zone.Now()

			// collect request data
			requestHeader := c.Request().Header
			requestData, err := c.GetRawData()
			if err != nil {
				fmt.Println(err.Error())
				c.Next()
			}
			r := c.Request()
			r.Body = ioutil.NopCloser(bytes.NewBuffer(requestData)) // key point
			c.SetRequest(r)

			// collect response data
			responseWriter := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer()}
			c.SetWriter(responseWriter)

			// print request
			defer log.Info(c.ClientIP(), toto.V{
				"Method":         c.Request().Method,
				"Path":           c.Request().RequestURI,
				"Proto":          c.Request().Proto,
				"Status":         responseWriter.Status(),
				"UA":             c.Request().UserAgent(),
				"Latency":        zone.Now().Sub(startedAt),
				"RequestHeader":  requestHeader,
				"RequestBody":    string(requestData),
				"ResponseHeader": responseWriter.Header(),
				"ResponseBody":   responseWriter.body.String(),
			})
		}

		c.Next()

		// after request

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
