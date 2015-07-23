// module.go

package module

import (
	"agent"
	"log"
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
func PeriodicScheduler(cb interface{}, ag agent.Agent, seconds int16) {
	callback, ok := cb.(func(agent.Agent, net.Conn, map[string]interface{}))
	if ok == false {
		log.Fatal("In PeriodicScheduler: Fail to convert callback function")
	}

	timer := time.NewTicker(time.Duration(seconds) * time.Second)

	for {
		select {
		case <-timer.C:
			go callback(ag, nil, nil)
		}
	}
}
