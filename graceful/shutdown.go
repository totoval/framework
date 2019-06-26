package graceful

import (
	"github.com/totoval/framework/cache"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/m"
	"github.com/totoval/framework/logs"
	"github.com/totoval/framework/queue"
)

func ShutDown(quietly bool) {
	logInfo(quietly, "Totoval is shutting down")
	closeQueue(quietly)
	closeCache(quietly)
	closeDB(quietly)
	logInfo(quietly, "Totoval is shut down")
}

func closeQueue(quietly bool) {
	defer panicRecover(quietly)
	logInfo(quietly, "Queue closing")
	if err := queue.Queue().Close(); err != nil {
		logFatal(quietly, "queue close failed", logs.Field{"error": err})
	}
	logInfo(quietly, "Queue closed")
}
func closeDB(quietly bool) {
	defer panicRecover(quietly)
	logInfo(quietly, "Database closing")
	if err := m.H().DB().Close(); err != nil {
		logFatal(quietly, "database close failed", logs.Field{"error": err})
	}
	logInfo(quietly, "Database closed")
}
func closeCache(quietly bool) {
	defer panicRecover(quietly)
	logInfo(quietly, "Cache closing")
	if err := cache.Cache().Close(); err != nil {
		logFatal(quietly, "cache close failed", logs.Field{"error": err})
	}
	logInfo(quietly, "Cache closed")
}

func panicRecover(quietly bool) {
	if err := recover(); err != nil {
		logFatal(quietly, "Totoval shutting down failed", logs.Field{"error": err})
	}
}

func logInfo(quietly bool, msg string, v ...logs.Field) {
	if !quietly {
		log.Info(msg, v...)
	}
}
func logFatal(quietly bool, msg string, v ...logs.Field) {
	if !quietly {
		log.Fatal(msg, v...)
	}
}
