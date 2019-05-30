package cache

import (
	"github.com/totoval/framework/cache/driver/memory"
	"github.com/totoval/framework/cache/driver/redis"

	"github.com/totoval/framework/config"
)

var cer cacher

func Initialize() {
	cer = setStore("default")
}
func setStore(store string) (cer cacher) {

	_conn := store
	if store == "default" {
		_conn = config.GetString("cache." + store)
		if _conn == "" {
			panic("cache connection parse error")
		}
	}

	// get driver instance and connect cache store
	switch _conn {
	case "memory":
		cer = memory.NewMemory(
			config.GetString("cache.stores.memory.prefix"),
			config.GetUint("cache.stores.memory.default_expiration_minute"),
			config.GetUint("cache.stores.memory.cleanup_interval_minute"),
		)
	case "redis":
		connection := config.GetString("cache.stores.redis.connection") // cache
		cer = redis.NewRedis(
			config.GetString("database.redis."+connection+".host"),
			config.GetString("database.redis."+connection+".port"),
			config.GetString("database.redis."+connection+".password"),
			config.GetInt("database.redis."+connection+".database"),
			config.GetString("database.redis.options.prefix"),
		)
	default:
		panic("incorrect cache connection provided")
	}

	return cer
}

func Store(store string) cacher {
	return setStore(store)
}

func Cache() cacher {
	return cer
}
