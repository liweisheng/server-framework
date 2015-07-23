package agent

import (
	"fmt"
	"net"
	"os"
	"strconv"
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

func (ma *MonitorAgent) IsMaster() bool {
	return false
}

/// 获得加载当前monitor组建的服务的信息,由于monitor组建同时被master加载,
/// 所以对于 master服务器其返回的信息就是master的信息.
func (ma *MonitorAgent) GetServerInfo() map[string]interface{} {
	return ma.ServerInfo
}

/// 获得master服务器的信息.
func (ma *MonitorAgent) GetMasterInfo() map[string]interface{} {
	return ma.MasterInfo
}

/// 给master发送通知消息,发送成功后返回.
///
/// @param moduleID 对应master端处理该通知的moduleID
/// @param msg 通知的消息体.
func (ma *MonitorAgent) Notify(moduleID string, msg map[string]string) {}

/// 发送请求给master,请求发送成功后并不等待到master端回复，而是直接返回.
///
/// @param moduleID 是请求的master端的module id
/// @param msg 请求的消息
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

/// 连接Master
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
	hostAndPort := host_s + ":" + strconv.Itoa(port_i)
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
