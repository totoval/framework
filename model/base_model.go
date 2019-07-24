package model

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/jinzhu/gorm"

	"github.com/iancoleman/strcase"

	"github.com/totoval/framework/database"
)

type BaseModeller interface {
	DB() *gorm.DB
	SetTX(db *gorm.DB)
}

type BaseModel struct {
	db *gorm.DB
}

func (bm *BaseModel) DB() *gorm.DB {
	if bm.db == nil {
		return database.DB()
	}
	return bm.db
}
func (bm *BaseModel) SetTX(db *gorm.DB) {
	bm.SetDB(db)
}
func (bm *BaseModel) SetDB(db *gorm.DB) {
	bm.db = db
}
func (bm *BaseModel) SetTableName(tableName string) string {
	return fmt.Sprintf("%s%s", database.Prefix(), tableName)
}
func (bm *BaseModel) BeforeSave(scope *gorm.Scope) (err error) {
	defer func() {
		if _err := recover(); _err != nil {
			if __err, ok := _err.(error); ok {
				err = __err
				return
			}
			err = errors.New(fmt.Sprint(_err))
			return
		}
	}()

	callMutator(scope, false)
	return nil
}
func (bm *BaseModel) BeforeUpdate(scope *gorm.Scope) (err error) {
	defer func() {
		if _err := recover(); _err != nil {
			if __err, ok := _err.(error); ok {
				err = __err
				return
			}
			err = errors.New(fmt.Sprint(_err))
			return
		}
	}()

	callMutator(scope, false)
	return nil
}

//func (bm *BaseModel) AfterUpdate() error {
//
//}
//func (bm *BaseModel) AfterSave() error {
//
//}
func (bm *BaseModel) AfterFind(scope *gorm.Scope) (err error) {
	defer func() {
		if _err := recover(); _err != nil {
			if __err, ok := _err.(error); ok {
				err = __err
				return
			}
			err = errors.New(fmt.Sprint(_err))
			return
		}
	}()

	callMutator(scope, true)
	return nil
}

func callMutator(scope *gorm.Scope, isGetter bool) {

	reflectValue := scope.IndirectValue()
	// Only get address from non-pointer
	if reflectValue.CanAddr() && reflectValue.Kind() != reflect.Ptr {
		reflectValue = reflectValue.Addr()
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < reflectValue.Type().Elem().NumField(); i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, index int) {
			fieldName := reflectValue.Type().Elem().Field(index).Name
			fieldValue := reflectValue.Elem().Field(index).Interface()

			if isGetter {
				getter(&reflectValue, fieldName, fieldValue)
			} else {
				setter(reflectValue, fieldName, fieldValue)
			}
			wg.Done()
		}(wg, i)
	}
	wg.Wait()
}

func setter(reflectValue reflect.Value, fieldName string, fieldValue interface{}) {

	const setterMethodTemplate = "Set%sAttribute"
	methodName := fmt.Sprintf(setterMethodTemplate, strcase.ToCamel(fieldName))

	if methodValue := reflectValue.MethodByName(methodName); methodValue.IsValid() {
		methodValue.Interface().(func(value interface{}))(fieldValue)
	}
}
func getter(reflectValue *reflect.Value, fieldName string, fieldValue interface{}) {

	const getterMethodTemplate = "Get%sAttribute"
	methodName := fmt.Sprintf(getterMethodTemplate, strcase.ToCamel(fieldName))

	if methodValue := reflectValue.MethodByName(methodName); methodValue.IsValid() {
		newGetData := methodValue.Interface().(func(value interface{}) interface{})(fieldValue)
		reflectValue.Elem().FieldByName(fieldName).Set(reflect.ValueOf(newGetData))
	}
}
