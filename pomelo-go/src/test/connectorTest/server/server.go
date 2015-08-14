package main

import (
	"connector"
	"connector/tcp_connector"
	"context"
	"fmt"
	"github.com/cihub/seelog"
	"os"
	"strconv"
)

/// 新消息到达时回调.
func NewMsgCB(id uint32, msg map[string]interface{}) {
	fmt.Printf("from sid<%v> decoded message:<%v>\n", id, msg)
}

/// 新连接到达时回调.
func NewConnCB(sock connector.Socket) {
	fmt.Printf("new connection<id:%v remoteAddr:%v>\n", sock.ID(), sock.RemoteAddress())
}

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "too few args,args form: <host port>\n")
		os.Exit(1)
	}
	host := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid port,need integer type,your input port: <port>\n", os.Args[2])
		os.Exit(1)
	}
	ctx := context.GetContext()
	currentServer := make(map[string]interface{})
	currentServer["id"] = "connector-1"
	currentServer["serverType"] = "connector"
	currentServer["host"] = "127.0.0.1"
	currentServer["port"] = 8888
	ctx.CurrentServer = currentServer
	defer seelog.Flush()

	tcp_cnct := tcp_connector.NewTcpConnector(host, port, nil)

	tcp_cnct.RegistNewConnCB(NewConnCB)
	tcp_cnct.RegistNewMsgCB(NewMsgCB)
	tcp_cnct.Start()
	ch := make(chan int)
	<-ch
}
