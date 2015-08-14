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
	allOpts := make(map[string]map[string]interface{})
	cosessOpts := make(map[string]interface{})

	//指定同一个uid可以绑定多个session
	cosessOpts["multiBind"] = "yes"
	allOpts["cosession"] = cosessOpts
	ctx.AllOpts = allOpts
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

	// 启动cosession，cosessionrpcserver,corpcserver
	cosess.Start()
	cosessRS.Start()
	corpcS.Start()

	cosess.CreateSession(1, "connector-1", nil)
	cosess.BindUID("Zhang San", 1)
	cosess.CreateSession(2, "connector-1", nil)
	cosess.BindUID("Li Si", 2)
	cosess.CreateSession(3, "connector-1", nil)
	cosess.BindUID("Li Si", 3)
	cosess.CreateSession(4, "connector-1", nil)
	cosess.BindUID("Wang Wu", 4)

	ch := make(chan int)
	<-ch
}
