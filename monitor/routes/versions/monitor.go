package versions

import (
	"github.com/totoval/framework/config"
	"github.com/totoval/framework/http/middleware"
	"github.com/totoval/framework/monitor/routes/groups"

	"github.com/totoval/framework/request"
	"github.com/totoval/framework/route"
)

func NewMonitor(engine *request.Engine) {
	ver := route.NewVersion(engine, "monitor")

	accounts := make(map[string]string)
	accounts[config.GetString("monitor.username")] = config.GetString("monitor.password")

	// noauth routes
	ver.NoAuth("", func(grp route.Grouper) {
		grp.AddGroup("/dashboard", &groups.DashboardGroup{})
	}, middleware.BasicAuth(accounts))
}
