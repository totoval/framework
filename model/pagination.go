package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
	"math"
	"reflect"
	"strconv"
)

type Model gorm.DB

type Pagination struct {
	currentPageItemCount uint
	currentPageNum       uint
	totalPageNum         uint
	totalItemCount       uint
	itemArr              interface{}
	perPage              uint
}

func (p *Pagination) Count() uint {
	return p.currentPageItemCount
}
func (p *Pagination) CurrentPage() uint {
	return p.currentPageNum
}
func (p *Pagination) LastPage() uint {
	return p.totalPageNum
}
func (p *Pagination) FirstItem() interface{} {
	if len(p.itemArr.([]interface{})) > 0 {
		return p.itemArr.([]interface{})[0]
	}
	return nil
}
func (p *Pagination) LastItem() interface{} {
	if len(p.itemArr.([]interface{})) > 0 {
		return p.itemArr.([]interface{})[len(p.itemArr.([]interface{}))-1]
	}
	return nil
}
func (p *Pagination) ItemArr() interface{} {
	return p.itemArr
}
func (p *Pagination) Total() uint {
	return p.totalItemCount
}
func (p *Pagination) PerPage() uint {
	return p.perPage
}

// Model(*Q(&User{}, data, []Sort{}, 1, false)).Paginate(c, perPage)
func (bm *Model) Paginate(model interface{}, c *gin.Context, perPage uint) (pagination Pagination, err error) {
	// validate paginate params
	p := c.DefaultQuery("page", "1")
	if err := validator.New().Var(p, "numeric,gt=0"); err != nil {
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
		return pagination, err
	}

	_p, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		return pagination, err
	}
	page := uint(_p)

	// perpage
	pagination.perPage = perPage

	// current page num
	pagination.currentPageNum = page

	// calc total item count
	count := gorm.DB(*bm)
	if err = count.Model(model).Count(&pagination.totalItemCount).Error; err != nil {
		return pagination, err
	}

	// calc total page num
	pagination.totalPageNum = uint(math.Ceil(float64(pagination.totalItemCount) / float64(perPage)))

	// get data
	data := gorm.DB(*bm)
	if err = data.Offset(perPage * (page - 1)).Limit(perPage).Find(model).Error; err != nil {
		return pagination, err
	}
	pagination.itemArr = model

	// get currentPageItemCount
	if reflect.ValueOf(model).Elem().Type().Kind() != reflect.Slice {
		return pagination, errors.New("result is not a slice")
	}
	s := reflect.ValueOf(model).Elem()
	pagination.currentPageItemCount = uint(s.Len())

	return pagination, nil
}

//func (m *Model) shouldInstantiate() { //     private function shouldInstantiate(bool $should, $primary_key_variable = null)
//
//}
//
//func (m *Model) readOnlyGuardian() { //     private function readOnlyGuardian()
//
//}
