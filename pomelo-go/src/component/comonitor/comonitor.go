package comonitor

import (
	"context"
	"github.com/cihub/seelog"
	"pomelo_admin"
)

type CoMonitor struct {
	*pomelo_admin.MonitorConsoleService
}

func NewCoMonitor() *CoMonitor {
	ctx := context.GetContext()
	coMonitor, ok := ctx.GetComponent("comonitor").(*CoMonitor)
	if ok == true {
		return coMonitor
	}

	mcs := pomelo_admin.NewMonitorConsoleService(ctx)
	coMonitor = &CoMonitor{mcs}
	seelog.Infof("<%v> component CoMonitor created...", ctx.GetServerID())
	ctx.RegisteComponent("comonitor", coMonitor)
	return coMonitor
}
