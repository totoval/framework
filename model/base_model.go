package model

import (
	"errors"
	"gopkg.in/go-playground/validator.v9"
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

		// fill data
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

		// fill default
		for j := 0; j < fillType.NumField(); j++ {
			if newDataType.Field(i).Type == fillType.Field(j).Type &&newDataType.Field(i).Name == fillType.Field(j).Name {
				newDataValue.Field(i).Set(fillValue.Field(j))
				break
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

	out = &inData
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
	if err := db.Save(&inData).Error; err != nil {
		return err
	}

	out = inData
	return nil
}

// out must be a struct pointer
func First(out interface{}) (error) {
	if err := db.Where(out).First(out).Error; err != nil {
		return err
	}
	return nil
}


func (m *BaseModel) DeleteObj() { // public function deleteObj(int $id = null, bool $force_delete = false): bool

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
