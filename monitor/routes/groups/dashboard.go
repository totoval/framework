package groups

import (
	"github.com/totoval/framework/monitor/app/http/controllers"

	"github.com/totoval/framework/route"
)

type DashboardGroup struct {
	DashboardController          controllers.Dashboard
	DashboardWebsocketController controllers.DashboardWebsocketController
}

func (dg *DashboardGroup) Group(group route.Grouper) {
	group.GET("/", dg.DashboardController.Index)
	group.Websocket("/ws", &dg.DashboardWebsocketController)
}
