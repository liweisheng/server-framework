// module.go

package module

import (
	"agent"
	"net"
	"time"
)

type Module interface {
	MonitorHandler(agent agent.Agent, conn net.Conn, msg map[string]interface{})
	MasterHandler(agent agent.Agent, conn net.Conn, msg map[string]interface{})
	ClientHandler(agent agent.Agent, conn net.Conn, msg map[string]interface{})
	Start()
	ModuleID() string
	GetType() string
	GetInterval() int16
}

///以seconds为周期调用cb
func PeriodicScheduler(cb interface{}, ag interface{}, seconds int16) {
	callback := cb.(func(interface{}, net.Conn, map[string]interface{}))
	ag_ := ag.(*agent.Agent)

	timer := time.NewTicker(time.Duration(seconds) * time.Second)

	for {
		select {
		case <-timer.C:
			go callback(ag_, nil, nil)
		}
	}
}
