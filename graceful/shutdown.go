package graceful

import (
	"context"

	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/m"
	"github.com/totoval/framework/logs"
	"github.com/totoval/framework/queue"
)

func ShutDown(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			closeQueue()
			closeDB()
			closeCache()
			return
		default:
			log.Info("Totoval Running")
		}
	}
}

func closeQueue() {
	log.Info("Queue closing")
	if err := queue.Queue().Close(); err != nil {
		log.Fatal("queue close failed", logs.Field{"error": err})
	}
	log.Info("Queue closed")
}
func closeDB() {
	log.Info("Database closing")
	if err := m.H().DB().Close(); err != nil {
		log.Fatal("database close failed", logs.Field{"error": err})
	}
	log.Info("Database closed")
}
func closeCache() {
	log.Info("Cache closing")

	//@todo close cache
	
	log.Info("Cache closed")
}
