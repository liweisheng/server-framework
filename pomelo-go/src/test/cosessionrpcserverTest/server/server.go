package main

import (
	"component/corpcserver"
	"component/cosession"
	"component/cosessionrpcserver"
	"context"
	"github.com/cihub/seelog"
)

func main() {
	ctx := context.GetContext()
	defer seelog.Flush()
	currentServer := make(map[string]interface{})
	currentServer["id"] = "connector-1"
	currentServer["serverType"] = "connector"
	currentServer["host"] = "127.0.0.1"
	currentServer["port"] = 8888
	ctx.CurrentServer = currentServer
	cosess := cosession.NewCoSession()
	cosessRS := cosessionrpcserver.NewCoSessionRpcServer()
	corpcS := corpcserver.NewCoRpcServer()
	cosess.Start()
	cosessRS.Start()
	corpcS.Start()
	cosess.CreateSession(1, "connector-1", nil)
	cosess.BindUID("liweisheng", 1)
	ch := make(chan int)
	<-ch
}
