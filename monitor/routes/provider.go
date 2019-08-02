package routes

import (
	"github.com/totoval/framework/monitor/routes/versions"
	"github.com/totoval/framework/request"
	"github.com/totoval/framework/route"
)

func Register(router *request.Engine) {
	defer route.Bind(router)

	versions.NewMonitor(router)
}
