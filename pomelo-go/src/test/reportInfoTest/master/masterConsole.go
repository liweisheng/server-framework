package main

import (
	"context"
	// "module/reportInfo"
	// "os"
	"pomelo_admin"
	// "time"
)

func main() {
	ctx := context.NewContext()
	ctx.MasterInfo["id"] = "master"
	ctx.MasterInfo["host"] = "127.0.0.1"
	ctx.MasterInfo["port"] = 8888

	ctx.CurrentServer["id"] = "master-1"
	ctx.CurrentServer["host"] = "192.168.1.2"
	ctx.CurrentServer["port"] = 8889
	ctx.CurrentServer["clientPort"] = 8889
	ctx.CurrentServer["frontend"] = "false"
	ch := make(chan int)

	mcs := pomelo_admin.NewMasterConsoleService(ctx)
	mcs.Start()
	select {
	case <-ch:
	}
}
