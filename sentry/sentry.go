package sentry

import (
	"fmt"

	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/config"
)

func Initialize() {
	if config.GetBool("sentry.enable") {
		err := raven.SetDSN(fmt.Sprintf("https://%s:%s@%s/%s",
			config.GetString("sentry.key"),
			config.GetString("sentry.secret"),
			config.GetString("sentry.host"),
			config.GetString("sentry.project"),
		))
		if err != nil {
			panic(err)
		}
	}
}

func Use(r *gin.Engine, onlySendOnCrash bool) {
	if config.GetBool("sentry.enable") {
		r.Use(sentry.Recovery(raven.DefaultClient, onlySendOnCrash))
	}
}

func CaptureError(err error) {
	if config.GetBool("sentry.enable") {
		raven.CaptureErrorAndWait(err, map[string]string{
			"env": config.GetString("app.env"),
		})
	}
}

func CapturePanic(handler func()) {
	if config.GetBool("sentry.enable") {
		raven.CapturePanic(handler, map[string]string{
			"env": config.GetString("app.env"),
		})
	}
}
