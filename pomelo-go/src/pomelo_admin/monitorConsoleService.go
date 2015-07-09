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

func NewMonitorConsoleService(ctx *context.Context, servInfo map[string]interface{}) *MonitorConsoleService {
	moduleMap := make(map[string]module.Module)
	monitorAgent := agent.NewMonitorAgent(ctx.MasterInfo, servInfo)
	return &MonitorConsoleService{ctx, moduleMap, monitorAgent, SV_INIT}

}

func (m *MonitorConsoleService) handleConnection(conn net.Conn) {
	if m.status != SV_START {
		return
	}

	handlerConnectionRecv(m, conn)
}

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
