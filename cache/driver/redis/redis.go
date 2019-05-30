package redis

import (
	"errors"

	r "github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"

	. "github.com/totoval/framework/cache/utils"
	"github.com/totoval/framework/helpers/zone"
)

func NewRedis(host string, port string, password string, dbIndex int, prefix string) *redis {
	client := r.NewClient(&r.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       dbIndex,
	})

	//pong, err := client.Ping().Result()
	//fmt.Println(pong, err)
	//// Output: PONG <nil>

	return &redis{
		redisBasic{
			cache:  client,
			prefix: prefix,
		},
	}
}

type redis struct {
	redisBasic
}

func (re *redis) Pget(key string, valuePtr proto.Message, defaultValuePtr ...proto.Message) error {
	var valueBytes []byte

	k := NewKey(key, re.Prefix())

	if !re.Has(k.Raw()) {
		//@todo Event CacheMissed
		if len(defaultValuePtr) > 0 {
			return copier.Copy(valuePtr, defaultValuePtr[0])
		}
		return errors.New("key not exist")
	}
	err := re.cache.Get(k.Prefixed()).Scan(&valueBytes)
	if err != nil {
		return err
	}

	//@todo Event CacheHit
	if err := proto.Unmarshal(valueBytes, valuePtr); err != nil {
		return err
	}
	return nil
}

// ------------------------------------------------------------------------------
// the same
// ------------------------------------------------------------------------------

func (re *redis) Ppull(key string, valuePtr proto.Message, defaultValuePtr ...proto.Message) error {
	k := NewKey(key, re.Prefix())

	err := re.Pget(k.Raw(), valuePtr, defaultValuePtr...)
	if err != nil {
		return err
	}

	re.Forget(k.Raw())

	return nil
}
func (re *redis) Pput(key string, valuePtr proto.Message, future zone.Time) bool {
	valueBytes, err := proto.Marshal(valuePtr)
	if err != nil {
		return false
	}
	return re.Put(key, valueBytes, future)
}
func (re *redis) Padd(key string, valuePtr proto.Message, future zone.Time) bool {
	valueBytes, err := proto.Marshal(valuePtr)
	if err != nil {
		return false
	}
	return re.Add(key, valueBytes, future)
}
func (re *redis) Pforever(key string, valuePtr proto.Message) bool {
	valueBytes, err := proto.Marshal(valuePtr)
	if err != nil {
		return false
	}
	return re.Forever(key, valueBytes)
}
