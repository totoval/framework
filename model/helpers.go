package model

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"strings"
)

type sortDirection byte

func (sd sortDirection) String() string {
	switch sd {
	case ASC:
		return "asc"
	case DESC:
		return "desc"
	}
	panic(errors.New("type sortDirection parsed error"))
}

type Sort struct {
	Key       string
	Direction sortDirection
}

const (
	ASC sortDirection = iota
	DESC
)

func structToMap(data interface{}) (result map[string]interface{}) {
	result = map[string]interface{}{}

	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)

	for i := 0; i < dataType.NumField(); i++ {
		dataField := dataType.Field(i)
		dataFieldValue := dataValue.Field(i)

		switch dataFieldValue.Kind() {
		case reflect.Bool:
			result[dataField.Name] = dataFieldValue.Bool()
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result[dataField.Name] = dataFieldValue.Int()
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result[dataField.Name] = dataFieldValue.Uint()
			break
		case reflect.String:
			result[dataField.Name] = dataFieldValue.String()
			break
		case reflect.Slice: //TODO...
		case reflect.Map: //TODO...
		case reflect.Struct: //TODO...
		default:
			panic("struct value cannot be slice, map or struct")
		}
	}

	return
}

func fillStruct(data interface{}, fill interface{}, mustFill bool) (interface{}, error) {
	dataType := reflect.TypeOf(data).Elem()
	dataValue := reflect.ValueOf(data).Elem()
	fillType := reflect.TypeOf(fill)
	fillValue := reflect.ValueOf(fill)

	newDataType := reflect.TypeOf(fill)
	newDataValue := reflect.New(reflect.TypeOf(fill)).Elem()
	for i := 0; i < newDataType.NumField(); i++ {
		if !newDataValue.Field(i).IsValid() || !newDataValue.Field(i).CanSet() {
			return nil, errors.New("model value cannot be filled")
		}

		// fill original data
		isFilled := false
		for j := 0; j < dataType.NumField(); j++ {
			if newDataType.Field(i).Type == dataType.Field(j).Type && newDataType.Field(i).Name == dataType.Field(j).Name {
				newDataValue.Field(i).Set(dataValue.Field(j))
				isFilled = true
				break
			}
		}
		if !mustFill && isFilled {
			continue
		}

		// fill input data
		for j := 0; j < fillType.NumField(); j++ {

			// if kind is ptr
			if fillValue.Field(j).Kind() == reflect.Ptr {

				// if not set, then do not fill
				if fillValue.Field(j).IsNil() {
					continue
				}

				// fill data
				if newDataType.Field(i).Type == fillType.Field(j).Type && newDataType.Field(i).Name == fillType.Field(j).Name {
					newDataValue.Field(i).Set(fillValue.Field(j))
					break
				}
			} else {
				// if kind is value fill
				//@todo WARN we have not check zero value for update, so if the type is `string`, and its' value is `""`(not set @fill), we'll override the correct value!!!!
				// fill data
				if newDataType.Field(i).Type == fillType.Field(j).Type && newDataType.Field(i).Name == fillType.Field(j).Name {
					newDataValue.Field(i).Set(fillValue.Field(j))
					break
				}
			}

		}

	}

	return newDataValue.Addr().Interface(), nil
}

func Transaction(f func(), attempts uint) {
	if attempts <= 0 {
		attempts = 1
	}
	var currentAttempt uint
	currentAttempt = 1
	tx := db.Begin()
	defer func(tx *gorm.DB) {
		if err := recover(); err != nil {
			var __err error
			if _err, ok := err.(error); ok {
				__err = _err
			} else {
				__err = errors.New(err.(string)) //@todo err.(string) may be down when `panic(123)`
			}
			handleTransactionException(tx, f, __err, currentAttempt, attempts)
		}
	}(tx)
	f()
	tx.Commit()
}
func handleTransactionException(tx *gorm.DB, f func(), err error, currentAttempt uint, maxAttempts uint) {
	tx.Rollback()
	if currentAttempt < maxAttempts {
		Transaction(f, maxAttempts-currentAttempt)
	}

	panic(err)
}

// out must be a struct pointer
func Create(out interface{}) error {
	//dataMap := structToMap(data)

	// fill default data
	defaultData := out.(Modeller)
	inData, err := fillStruct(out, defaultData.Default(), false)
	if err != nil {
		return err
	}

	// validate data
	if err := validator.New().Struct(inData); err != nil {
		//// this check is only needed when your code could produce
		//// an invalid value for validation such as interface with nil
		//// value most including myself do not usually have code like this.
		//if _, ok := err.(*validator.InvalidValidationError); ok {
		//	fmt.Println(err)
		//	return
		//}
		//
		//for _, err := range err.(validator.ValidationErrors) {
		//
		//	fmt.Println(err.Namespace())
		//	fmt.Println(err.Field())
		//	fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
		//	fmt.Println(err.StructField())     // by passing alt name to ReportError like below
		//	fmt.Println(err.Tag())
		//	fmt.Println(err.ActualTag())
		//	fmt.Println(err.Kind())
		//	fmt.Println(err.Type())
		//	fmt.Println(err.Value())
		//	fmt.Println(err.Param())
		//	fmt.Println()
		//}
		//
		//// from here you can create your own error messages in whatever language you wish
		return err
	}

	// create record
	if err := db.Create(inData).Error; err != nil {
		return err
	}

	copier.Copy(out, inData)
	return nil
}

// out must be a struct pointer
func Save(out interface{}, modify interface{}) error {
	// modify data
	inData, err := fillStruct(out, modify, true)
	if err != nil {
		return err
	}

	// validate data
	if err := validator.New().Struct(inData); err != nil {
		//// this check is only needed when your code could produce
		//// an invalid value for validation such as interface with nil
		//// value most including myself do not usually have code like this.
		//if _, ok := err.(*validator.InvalidValidationError); ok {
		//	fmt.Println(err)
		//	return
		//}
		//
		//for _, err := range err.(validator.ValidationErrors) {
		//
		//	fmt.Println(err.Namespace())
		//	fmt.Println(err.Field())
		//	fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
		//	fmt.Println(err.StructField())     // by passing alt name to ReportError like below
		//	fmt.Println(err.Tag())
		//	fmt.Println(err.ActualTag())
		//	fmt.Println(err.Kind())
		//	fmt.Println(err.Type())
		//	fmt.Println(err.Value())
		//	fmt.Println(err.Param())
		//	fmt.Println()
		//}
		//
		//// from here you can create your own error messages in whatever language you wish
		return err
	}

	// save record
	if err := db.Where(out).Save(inData).Error; err != nil {
		return err
	}

	copier.Copy(out, inData)
	return nil
}

func SaveByID(id interface{}, out interface{}, modify interface{}) error {
	//@todo First(), get primarykey through tag, then save
	return Save(out, modify)
}

// out must be a struct pointer
func First(out interface{}, withTrashed bool) (error) {
	_db := db
	if withTrashed {
		_db = _db.Unscoped()
	}
	if err := _db.Where(out).First(out).Error; err != nil {
		return err
	}
	return nil
}

func Delete(in interface{}, force bool) (error) {
	_db := db
	if force {
		_db = _db.Unscoped()
	}
	if err := _db.Delete(in).Error; err != nil {
		return err
	}
	return nil
}

func deleteKeyName(in interface{}) (string, error) {
	dataType, ok := reflect.TypeOf(in).Elem().FieldByName("DeletedAt")
	if !ok {
		return "", errors.New("cannot get DeletedAt key name type")
	}

	gormTag := dataType.Tag.Get("gorm")
	for _, gormTagType := range strings.Split(gormTag, ";") {
		tmp := strings.Split(gormTagType, ":")
		key := tmp[0]
		value := tmp[1]
		if key == "column" {
			return value, nil
		}
	}
	return "", errors.New("cannot get DeletedAt key name")
}

func Restore(in interface{}) (error) {
	deleteKeyName, err := deleteKeyName(in)
	if err != nil {
		return err
	}
	if err := db.Unscoped().Model(in).Update(deleteKeyName, gorm.Expr("NULL")).Error; err != nil {
		return err
	}
	return nil
}

func Q(filterArr []Filter, sortArr []Sort, limit int, withTrashed bool) *gorm.DB {
	_db := mapFilter(db, filterArr)

	for _, value := range sortArr {
		_db = _db.Order(value.Key + " " + value.Direction.String())
	}

	if limit != 0 {
		_db = _db.Limit(limit)
	}

	if withTrashed {
		_db = _db.Unscoped()
	}

	return _db
}

type Filter struct {
	Key string
	Op string
	Val interface{}
}
func mapFilter(_db *gorm.DB, filterArr []Filter) *gorm.DB {
	for _, filter := range filterArr {

		//if len(filter) > 3 {
		//	// error too many params
		//	panic(errors.New("too many params for sql filter"))
		//}

		switch filter.Op {
		// xxx is_null
		case "is_null":
			_db = _db.Where(filter.Key + " is null")
			break

		// xxx is_not_null
		case "is_not_null":
			_db = _db.Where(filter.Key + " is not null")
			break

		// xxx in array
		case "in":
			f, ok := filter.Val.([]interface{})
			if !ok {
				// error cannot parse array conditions for in
				panic(errors.New("cannot parse array conditions for IN"))
			}
			_db = _db.Where(filter.Key+" in (?)", f)
			break

		// xxx not in array
		case "not_in":
			f, ok := filter.Val.([]interface{})
			if !ok {
				// error cannot parse array conditions for in
				panic(errors.New("cannot parse array conditions for NOT IN"))
			}
			_db = _db.Where(filter.Key+" not in (?)", f)
			break

		// xxx between array
		case "between":
			f, ok := filter.Val.([]interface{})
			if !ok || len(f) != 2 {
				panic(errors.New("cannot parse array conditions for BETWEEN"))
			}
			_db = _db.Where(filter.Key + " between ? and ?", f[0], f[1])
			break

		// xxx like><!= yyy
		default:
			_db = _db.Where(filter.Key + " " + filter.Op + " ?", filter.Val)
		}
	}

	return _db
}

func Count(in Modeller, filterArr []Filter, withTrashed bool) (count uint, err error) {
	err = Q(filterArr, []Sort{}, 0, withTrashed).Model(&in).Count(&count).Error
	return count, err
}

func Exist(in Modeller, withTrashed bool) (exist bool) {
	err := First(&in, withTrashed)
	if err == nil {
		return true
	}
	return false
}
