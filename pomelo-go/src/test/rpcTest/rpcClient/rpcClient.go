package main

import (
	"context"
	"fmt"
	seelog "github.com/cihub/seelog"
	"log"
	"rpcclient"
)

func main() {
	ctx := context.NewContext()

	currentServerInfo := make(map[string]interface{})

	currentServerInfo["host"] = "127.0.0.1"
	currentServerInfo["port"] = 8888

	ctx.CurrentServer = currentServerInfo
	client := rpcclient.NewRpcClient(ctx)

	seelog.Errorf("Error,this is just for test")

	var reply string
	callRst := client.RpcCall("self", "Echo.Hi", "Eric", &reply)

	if callRst == nil {
		log.Fatal("What Fuck")
	}
	<-callRst.Done

	if callRst.Error != nil {
		log.Fatal(callRst.Error.Error())
	}

	fmt.Printf("Reply is : %v\n", *callRst.Reply.(*string))
}
