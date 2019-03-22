package cache

import (
    "github.com/totoval/framework/cache/driver"

    "github.com/totoval/framework/config"
)

var cer cacher

func Initialize(){
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