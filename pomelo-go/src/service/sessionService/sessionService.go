/*
author:liweisheng date:2015/07/08
*/

package sessionService

import (
	"connector/tcp_connector"
)

const (
	ST_CLOSED = iota
	ST_INITED = iota
)

type Session struct {
	status     int8
	id         int32
	uid        string
	frontendID string
	socket     *tcp_connector.TcpSocket
}

func NewSession(id int32, uid string, frontendId string, socket *tcp_connector.TcpSocket){
	return &{ST_INITED,id,uid,frontendId,socket}
}
