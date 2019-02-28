package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math"
	"reflect"
	"strconv"
)

type Model gorm.DB

type Pagination struct {
	currentPageItemCount uint
	currentPageNum uint
	totalPageNum uint
	totalItemCount uint
	itemArr interface{}
	perPage uint
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
	if len(p.itemArr.([]interface{})) > 0{
		return p.itemArr.([]interface{})[0]
	}
	return nil
}
func (p *Pagination) LastItem() interface{} {
	if len(p.itemArr.([]interface{})) > 0{
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
	type Page struct {
		Page uint `form:"page" json:"page" binding:"numeric,gt=0"`
	}
	var paginate Page
	//if err := c.BindQuery(&paginate); err != nil {
	//	return pagination, err
	//}

	page, err := strconv.ParseUint(c.Query("page"), 10, 32)
	if err != nil{
		return pagination, err
	}
	paginate.Page = uint(page)

	// perpage
	pagination.perPage = perPage

	// current page num
	pagination.currentPageNum = paginate.Page

	//// calc total item count
	//count := gorm.DB(*bm)
	//if err = count.Count(&out).Error; err != nil{
	//	return pagination, err
	//}

	// calc total page num
	pagination.totalPageNum = uint(math.Ceil(float64(pagination.totalItemCount) / float64(perPage)))

	// get data
	data := gorm.DB(*bm)
	if err = data.Offset(perPage * paginate.Page).Limit(perPage).Find(model).Error; err != nil{
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

type Modeller interface {
	Default() interface{}
	ObjArr(filterArr []Filter, sortArr []Sort, limit int, withTrashed bool) (interface{}, error)
	ObjArrPaginate(c *gin.Context, perPage uint, filterArr []Filter, sortArr []Sort, limit int, withTrashed bool) (pagination Pagination, err error)
}


//func (m *Model) shouldInstantiate() { //     private function shouldInstantiate(bool $should, $primary_key_variable = null)
//
//}
//
//func (m *Model) readOnlyGuardian() { //     private function readOnlyGuardian()
//
//}
