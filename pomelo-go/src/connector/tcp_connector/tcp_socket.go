/*
author:liweisheng date:2015/07/08
*/

/*
实现tcp socket相关的方法
*/

package tcp_connector

import (
	"connector"
	"context"
	"net"
)

type TcpSocket struct {
	socket     *net.TCPConn
	id         int32
	remoteAddr map[string]interface{}
	status     int8
}

//创建新的TcpSocket.参数id,sock不可以为空
func NewTcpSocket(id int32, sock *net.TCPConn) *TcpSocket {
	remoteAddr := make(map[string]interface{})
	addr := sock.RemoteAddr()
	host, port, err := net.SplitHostPort(addr.String())
	context.CheckError(err)

	remoteAddr["host"] = host
	remoteAddr["port"] = port

	return &TcpSocket{sock, id, remoteAddr, connector.ST_INITED}
}

func (ts *TcpSocket) Socket() interface{} {
	return ts.socket
}

func (ts *TcpSocket) ID() int32 {
	return ts.id
}

func (ts *TcpSocket) RemoteAddr() map[string]interface{} {
	return ts.remoteAddr
}

func (ts *TcpSocket) Send(msg []byte) (int, error) {
	if ts.status != connector.ST_INITED {
		return -1, nil
	}

	return ts.socket.Write(msg)
}

func (ts *TcpSocket) SendBatch(msgs ...[]byte) {}
func (ts *TcpSocket) Receive(recv []byte) (int, error) {
	if ts.status != connector.ST_INITED {
		return -1, nil
	}

	return ts.socket.Read(recv)
}

func (ts *TcpSocket) Disconnect() {
	if ts.status == connector.ST_CLOSED {
		return
	}

	ts.status = connector.ST_CLOSED
	ts.socket.Close()
}
