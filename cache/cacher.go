package cache

import (
	"github.com/totoval/framework/cache/driver"
)

type cacher interface {
	driver.ProtoCacher
	driver.BasicCacher
}
