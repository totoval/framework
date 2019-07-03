package auth

import (
	"reflect"

	"github.com/totoval/framework/config"
)

func NewUser() interface{} {

	ptr := reflect.New(reflect.TypeOf(config.GetInterface("auth.model_ptr"))).Elem()
	val := reflect.New(reflect.TypeOf(config.GetInterface("auth.model_ptr")).Elem())
	ptr.Set(val)
	return ptr.Interface()
}
