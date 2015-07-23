package pomelo_admin

import (
	"agent"
	"context"
	// "encoding/json"
	"fmt"
	"module"
	"net"
	"os"
)

/// MasterConsoleService实现consoleService接口
type MonitorConsoleService struct {
	Context      *context.Context
	ModuleMap    map[string]module.Module
	monitorAgent *agent.MonitorAgent
	status       int8
}

func NewMonitorConsoleService(ctx *context.Context) *MonitorConsoleService {
	moduleMap := make(map[string]module.Module)
	monitorAgent := agent.NewMonitorAgent(ctx.MasterInfo, ctx.CurrentServer)
	return &MonitorConsoleService{ctx, moduleMap, monitorAgent, SV_INIT}

}

func (m *MonitorConsoleService) handleConnection(conn net.Conn) {
	if m.status != SV_START {
		return
	}

	handlerConnectionRecv(m, conn)
}

/// 启动monitor.
///
/// 启动monitor时会首先拿到所有注册的module，并挂载到MonitorConsoleService上,
/// 然后对于设置类型为“push”的module定时调度其MonitoHandler方法, 为设置interval时
/// 默认定时调度时间5秒.
func (m *MonitorConsoleService) Start() {
	if m.status == SV_START {
		return
	}

	for _, v := range m.Context.Modules {
		m.ModuleMap[v.ModuleID()] = v
	}

	conn, err := m.monitorAgent.Connect()
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error:Fail connect to master host:%v post:%v ", m.Context.MasterInfo["host"], m.Context.MasterInfo["port"])
		os.Exit(1)
	}

	go m.handleConnection(conn)

	for _, mod := range m.ModuleMap {
		if mod != nil {
			if mod.GetType() == "push" {
				interval := mod.GetInterval()

				if interval == 0 {
					interval = 5
				}

				go module.PeriodicScheduler(mod.MonitorHandler, m.monitorAgent, interval)
			}
		}
	}

}

func (m *MonitorConsoleService) Stop() {
	if m.status != SV_START {
		return
	}

	m.status = SV_CLOSE

	m.monitorAgent.Close()
}

func (m *MonitorConsoleService) GetAgent() agent.Agent {
	return m.monitorAgent
}

func (m *MonitorConsoleService) GetModuleByID(id string) module.Module {
	return m.ModuleMap[id]
}
