package main

import (
	"context"
	// "module/reportInfo"
	// "os"
	// "fmt"
	"log"
	"module/reportInfo"
	"pomelo_admin"
	// "time"
)

func main() {
	ctx := context.NewContext()
	ctx.MasterInfo["id"] = "master"
	ctx.MasterInfo["host"] = "127.0.0.1"
	ctx.MasterInfo["port"] = 8888

	ctx.CurrentServer["serverType"] = "connector"
	ctx.CurrentServer["id"] = "connector-1"
	ctx.CurrentServer["host"] = "192.168.1.2"
	ctx.CurrentServer["port"] = 8889
	ctx.CurrentServer["clientPort"] = 8889
	ctx.CurrentServer["frontend"] = "false"

	ri, err := reportInfo.NewReportInfo("reportInfo", "push", 5, "127.0.0.1", 6379)

	if err != nil {
		log.Fatal(err.Error())
	}

	ctx.RegisteModule(ri)
	ch := make(chan int)

	mcs := pomelo_admin.NewMonitorConsoleService(ctx)
	mcs.Start()
	select {
	case <-ch:
	}
}
