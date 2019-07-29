package groups

import (
	"github.com/totoval/framework/monitor/app/http/controllers"

	"github.com/totoval/framework/route"
)

type DashboardGroup struct {
	DashboardController controllers.Dashboard
}

func (dg *DashboardGroup) Group(group route.Grouper) {
	group.GET("/", dg.DashboardController.Index)
}
