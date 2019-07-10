package model

import "github.com/totoval/framework/context"

type Modeller interface {
	BaseModeller
	TableName() string
	Default() interface{}
	ObjArr(filterArr []Filter, sortArr []Sort, limit int, withTrashed bool) (interface{}, error)
	ObjArrPaginate(c Context, perPage uint, filterArr []Filter, sortArr []Sort, limit int, withTrashed bool) (pagination Pagination, err error)
}

type Context interface {
	context.RequestParamContextor
}
