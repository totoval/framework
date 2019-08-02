package route

import (
	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/route"
)

func Url(routeName string, param ...toto.S) (url string, err error) {
	if len(param) > 0 {
		return route.RouteNameMap.Get(routeName, param[0])
	}
	return route.RouteNameMap.Get(routeName, toto.S{})
}
