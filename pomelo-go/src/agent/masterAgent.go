package agent

import (
	"fmt"
	"net"
	"sync"
)

type MasterAgent struct {
	Conns      map[net.Conn]bool
	MasterInfo map[string]interface{}

	Ch   chan int8
	lock *sync.RWMutex
}

func (ma *MasterAgent) NotifyAll() {
	for k, v := range ma.Conns {
		if v {
			_, err := k.Write([]byte("notify"))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func (ma *MasterAgent) GetMasterInfo() map[string]interface{} {
	return ma.MasterInfo
}

func (ma *MasterAgent) AddConnection(conn net.Conn) {
	ma.lock.Lock()
	defer ma.lock.Unlock()
	ma.Conns[conn] = true
}

func (ma *MasterAgent) RemoveConnection(conn net.Conn) {
	ma.lock.Lock()
	defer ma.lock.Unlock()
	delete(ma.Conns, conn)
}

func (ma *MasterAgent) IsMaster() bool {
	return true
}

func NewMasterAgent(ch chan int8, masterInfo map[string]interface{}) *MasterAgent {
	var conns = make(map[net.Conn]bool)
	return &MasterAgent{conns, masterInfo, ch, &sync.RWMutex{}}
}
