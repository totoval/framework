package model

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
	"strings"

	"reflect"
)

type BaseModel struct {}

type Modeller interface {
	Default() interface{}
	GetObjArr()         //@todo     public function getObjArr(?array $filter_arr = [], ?array $sort_arr = null, ?int $limit = null, bool $with_trashed = false): Collection;
	GetObjArrPaginate() //@todo     public function getObjArrPaginate(int $per_page, ?array $filter_arr = [], ?array $sort_arr = null, bool $with_trashed = false): LengthAwarePaginator;
}
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

func FillStruct(data interface{}, fill interface{}, mustFill bool) (interface{}, error) {
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
			}else{
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

// out must be a struct pointer
func Create(out interface{}) error {
	//dataMap := structToMap(data)

	// fill default data
	defaultData := out.(Modeller)
	inData, err := FillStruct(out, defaultData.Default(), false)
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
	inData, err := FillStruct(out, modify, true)
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

func (m *BaseModel) RestoreObj() { //     public function restoreObj(int $id = null): bool

}

func (m *BaseModel) CountObjArr() { //     public function countObjArr(?array $filter_arr = [], bool $with_trashed = false): int

}

func (m *BaseModel) DoFilterSortLimit() { //     protected function doFilterSortLimit(?array $filter_arr = [], ?array $sort_arr = null, ?int $limit = null, bool $with_trashed = false)

}

func (m *BaseModel) IsExistObjByID() { //     public function isExistObjByID(int $id, bool $with_trashed = false): bool

}
func (m *BaseModel) shouldInstantiate() { //     private function shouldInstantiate(bool $should, $primary_key_variable = null)

}

func (m *BaseModel) readOnlyGuardian() { //     private function readOnlyGuardian()

}
