package graceful

import (
	"github.com/totoval/framework/cache"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/m"
	"github.com/totoval/framework/logs"
	"github.com/totoval/framework/queue"
)

func ShutDown() {
	log.Info("Totoval is shutting down")
	closeQueue()
	closeCache()
	closeDB()
	log.Info("Totoval is shut down")
}

func closeQueue() {
	defer panicRecover()
	log.Info("Queue closing")
	if err := queue.Queue().Close(); err != nil {
		log.Fatal("queue close failed", logs.Field{"error": err})
	}
	log.Info("Queue closed")
}
func closeDB() {
	defer panicRecover()
	log.Info("Database closing")
	if err := m.H().DB().Close(); err != nil {
		log.Fatal("database close failed", logs.Field{"error": err})
	}
	log.Info("Database closed")
}
func closeCache() {
	defer panicRecover()
	log.Info("Cache closing")
	if err := cache.Cache().Close(); err != nil {
		log.Fatal("cache close failed", logs.Field{"error": err})
	}
	log.Info("Cache closed")
}

func panicRecover() {
	if err := recover(); err != nil {
		log.Fatal("Totoval shutting down failed", logs.Field{"error": err})
	}
}
