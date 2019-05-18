package driver

import (
	"github.com/totoval/framework/helpers/zone"
	"log"
	"strconv"
	"testing"
	"time"
)

// config
var hostname string = "127.0.0.1"
var port uint = 6379
var db uint = 0
var auth *string = nil
var prefix = "TEST_"

// test
var testKey string = "totoval_key"
var testString string = "hello totoval"
var testInt64 int64 = 8080
var testUint uint = 8080

func init() {
	zone.Initialize()
	c, err := NewRedis(hostname, port, db, auth, prefix, 20, 5)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	c.Forget(testKey)
	c.Close()
}

func TestString(t *testing.T) {
	zone.Initialize()
	c, err := NewRedis(hostname, port, db, auth, prefix, 20, 5)
	if err != nil {
		t.Error(err)
	}

	c.Put(testKey, testString, zone.Now().Add(zone.Second*60))
	val := c.Get(testKey)
	if val == nil {
		t.Error(val, "\t", testString)
	}
	if val.(string) != testString {
		t.Error(val.(string), "\t", testString)
	}
	c.Close()
}

func TestHas(t *testing.T) {
	zone.Initialize()
	c, err := NewRedis(hostname, port, db, auth, prefix, 20, 5)
	if err != nil {
		t.Error(err)
	}

	c.Put(testKey, testString, zone.Now().Add(zone.Second*60))
	val := c.Get(testKey)
	if val == nil {
		t.Error(val, "\t", testString)
	}
	if !c.Has(testKey) {
		t.Error(val, testKey)
	}
	c.Close()
}

func TestInt(t *testing.T) {
	zone.Initialize()
	c, err := NewRedis(hostname, port, db, auth, prefix, 20, 5)
	if err != nil {
		t.Error(err)
	}

	c.Put(testKey, testInt64, zone.Now().Add(zone.Second*60))
	val := c.Get(testKey)
	if val == nil {
		t.Error(val, "\t", testInt64)
	}
	if val.(string) != strconv.FormatUint(uint64(testInt64), 10) {
		t.Error(val.(string), "\t", testInt64)
	}
	c.Close()
}

func TestIntIncr(t *testing.T) {
	zone.Initialize()
	var num int64 = 5
	c, err := NewRedis(hostname, port, db, auth, prefix, 20, 5)
	if err != nil {
		t.Error(err)
	}

	c.Put(testKey, testInt64, zone.Now().Add(zone.Second*60))
	// incr
	c.Increment(testKey, num)

	// get key
	val := c.Get(testKey)
	if val == nil {
		t.Error(val, "\t", testInt64)
	}

	if err != nil {
		t.Error(err)
	}
	if val.(string) != strconv.FormatUint(uint64(testInt64+num), 10) {
		t.Error(val.(string), "\t", testInt64+num)
	}
	c.Close()
}

func TestIntDecr(t *testing.T) {
	zone.Initialize()
	var num int64 = 5
	c, err := NewRedis(hostname, port, db, auth, prefix, 20, 5)
	if err != nil {
		t.Error(err)
	}

	c.Put(testKey, testInt64, zone.Now().Add(zone.Second*60))
	// incr
	c.Decrement(testKey, num)

	// get key
	val := c.Get(testKey)
	if val == nil {
		t.Error(val, "\t", testInt64)
	}

	if err != nil {
		t.Error(err)
	}
	if val.(string) != strconv.FormatUint(uint64(testInt64-num), 10) {
		t.Error(val.(string), "\t", testInt64+num)
	}
	c.Close()
}

func TestPull(t *testing.T) {
	zone.Initialize()
	c, err := NewRedis(hostname, port, db, auth, prefix, 20, 5)
	if err != nil {
		t.Error(err)
	}

	c.Put(testKey, testInt64, zone.Now().Add(zone.Second*60))
	// incr
	c.Pull(testKey)

	// get key
	val := c.Get(testKey)
	if val != nil {
		t.Error(val)
	}
	c.Close()
}

func TestExpire(t *testing.T) {
	zone.Initialize()
	c, err := NewRedis(hostname, port, db, auth, prefix, 20, 5)
	if err != nil {
		t.Error(err)
	}

	c.Put(testKey, testInt64, zone.Now().Add(zone.Second*5))

	// slepp tow second
	time.Sleep(time.Second * 2)
	// get key
	val := c.Get(testKey)
	if val == nil {
		t.Error(val)
	}

	// slepp six second
	time.Sleep(time.Second * 5)
	val = c.Get(testKey)
	if val != nil {
		t.Error(val)
	}
	c.Close()
}
