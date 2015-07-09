package module

import (
	"agent"
	"net"
	// "time"
)

type Timer struct {
	ModuleId string
	Type     string
	Interval int16
}

func (t *Timer) MonitorHandler(agent agent.Agent, conn net.Conn, msg map[string]interface{}) {

}
func (t *Timer) MasterHandler(agent agent.Agent, conn net.Conn, msg map[string]interface{}) {}

func (t *Timer) ClientHandler(agent agent.Agent, conn net.Conn, msg map[string]interface{}) {

}

func (t *Timer) Start() {}
func (t *Timer) ModuleID() string {

	return t.ModuleId
}
func (t *Timer) GetType() string {
	return t.Type
}

func (t *Timer) GetInterval() int16 {
	return t.Interval
}
