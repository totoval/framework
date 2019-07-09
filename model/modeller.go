package model

import "github.com/totoval/framework/request"

type Modeller interface {
	BaseModeller
	TableName() string
	Default() interface{}
	ObjArr(filterArr []Filter, sortArr []Sort, limit int, withTrashed bool) (interface{}, error)
	ObjArrPaginate(c *request.Context, perPage uint, filterArr []Filter, sortArr []Sort, limit int, withTrashed bool) (pagination Pagination, err error)
}
