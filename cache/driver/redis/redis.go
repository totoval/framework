package redis

import (
	"errors"

	r "github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"

	"github.com/totoval/framework/cache/driver/common"
	. "github.com/totoval/framework/cache/utils"
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
		redisBasic: redisBasic{
			cache:  client,
			prefix: prefix,
		},
	}
}

type redis struct {
	redisBasic
	common.Common
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
