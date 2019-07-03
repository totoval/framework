package auth

import (
	"reflect"

	"github.com/totoval/framework/config"
)

func NewUser() interface{} {
	typeof := reflect.TypeOf(config.GetInterface("auth.model_ptr"))
	ptr := reflect.New(typeof).Elem()
	val := reflect.New(typeof.Elem())
	ptr.Set(val)
	return ptr.Interface()
}
