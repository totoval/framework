package config

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var v *viper.Viper

func init() {
	v = viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("json")
	v.AddConfigPath(".")
	err := v.ReadInConfig() // Find and read the config file
	if err != nil {         // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	v.AutomaticEnv()
}

func Env(envName string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return Get(envName, defaultValue[0])
	}
	return Get(envName)
}

func Add(name string, configuration map[string]interface{}) {
	v.Set(name, configuration)
}

func Get(path string, defaultValue ...interface{}) interface{} {
	if !v.IsSet(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return v.Get(path)
}
func GetInterface(path string, defaultValue ...interface{}) (value interface{}) {

	if len(defaultValue) > 0 {
		value = Get(path, defaultValue[0])
	} else {
		value = Get(path)
	}

	return value
}
func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(GetInterface(path, defaultValue...))
}
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(GetInterface(path, defaultValue...))
}
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(GetInterface(path, defaultValue...))
}
func GetUint(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(GetInterface(path, defaultValue...))
}
func GetBool(path string, defaultValue ...interface{}) bool {

	var value interface{}
	if len(defaultValue) > 0 {
		value = Get(path, defaultValue[0])
	} else {
		value = Get(path)
	}

	return cast.ToBool(value)
}
