package agent

import (
	"fmt"
	"net"
	"os"
)

type MonitorAgent struct {
	MasterInfo map[string]interface{}
	ServerInfo map[string]interface{}
	conn       net.Conn
	status     int8
}

func NewMonitorAgent(masterInfo map[string]interface{}, servInfo map[string]interface{}) *MonitorAgent {
	return &MonitorAgent{masterInfo, servInfo, nil, AG_INIT}
}

// func (ma *MonitorAgent) SetConn(conn net.Conn) {
// 	ma.conn = conn
// }

func (ma *MonitorAgent) Notify(moduleID string, msg map[string]string) {}

func (ma *MonitorAgent) Request(moduleID string, msg map[string]string) {}

func (ma *MonitorAgent) Close() {
	if ma.status != AG_START {
		return
	}

	ma.status = AG_CLOSE

	if ma.conn != nil {
		ma.conn.Close()
	}
}

func (ma *MonitorAgent) Connect() (net.Conn, error) {
	if ma.status == AG_START {
		return nil, nil
	}

	host, ok := ma.MasterInfo["host"]
	if !ok {
		fmt.Fprintf(os.Stdout, "Error: Master host not set\n")
		return nil, *new(error)
	}

	port, ok := ma.MasterInfo["port"]
	if !ok {
		fmt.Fprintf(os.Stdout, "Error: Master port not set\n")
		return nil, *new(error)
	}
	host_s := host.(string)
	port_i := port.(int)
	hostAndPort := host_s + ":" + string(port_i)
	masterAddr, err := net.ResolveTCPAddr("tcp4", hostAndPort)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: Resolve tcp4 addr error,error message:%v\n", err.Error())
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, masterAddr)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: Dial TCP error,error message:%v\n", err.Error())
		return nil, err
	}

	ma.conn = conn
	return conn, err

}
