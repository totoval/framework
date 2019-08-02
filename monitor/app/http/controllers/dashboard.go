package controllers

import (
	"net/http"

	"github.com/totoval/framework/config"
	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/http/controller"
	"github.com/totoval/framework/request"
)

type Dashboard struct {
	controller.BaseController
}

func (d *Dashboard) Index(c request.Context) {
	c.HTML(http.StatusOK, "totoval_dashboard.index", toto.V{
		"url": "ws://" + ":" + config.GetString("monitor.port") + "/monitor/dashboard/ws",
	})
	return
}
