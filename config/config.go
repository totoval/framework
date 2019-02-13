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
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	v.SetEnvPrefix("totoval")
	v.AutomaticEnv()
}

func Env(envName string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return Get(envName, defaultValue[0])
	}
	return Get(envName)
}

func Add(name string, configuration map[string]interface{}){
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

func GetString(path string, defaultValue ...interface{}) string {
	var value interface{}
	if len(defaultValue) > 0 {
		value = Get(path, defaultValue[0])
	}else{
		value = Get(path)
	}

	return cast.ToString(value)
}
func GetInt(path string, defaultValue ...interface{}) int {

	var value interface{}
	if len(defaultValue) > 0 {
		value = Get(path, defaultValue[0])
	}else{
		value = Get(path)
	}

	return cast.ToInt(value)
}