package cache

import (
	"github.com/totoval/framework/cache/driver"

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
		cer = driver.NewMemory(config.GetString("cache.stores.memory.prefix"), config.GetUint("cache.stores.memory.default_expiration_minute"), config.GetUint("cache.stores.memory.cleanup_interval_minute"))
		break
	case "redis":
		var err error
		cer, err = driver.NewRedis(config.GetString("cache.stores.redis.hostname"), config.GetUint("cache.stores.redis.port"), config.GetUint("cache.stores.redis.db"), config.GetString("cache.stores.redis.auth"), config.GetString("cache.stores.redis.prefix"), config.GetUint("cache.stores.redis.max_conn", 20), config.GetUint("cache.stores.redis.active_conn", 5))
		if err != nil {
			panic("cannot create or connect the redis driver")
		}
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
