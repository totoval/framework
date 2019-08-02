package monitor

import "github.com/totoval/framework/monitor/app/logics/dashboard"

func Shutdown() error {
	dashboard.Flow.Close()
	return nil
}
