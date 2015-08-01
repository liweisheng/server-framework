package corpcserver

import (
	"context"
	"github.com/cihub/seelog"
	"rpcserver"
)

type CoRpcServer struct {
	*rpcserver.RpcServer
}

func NewCoRpcServer() *CoRpcServer {
	ctx := context.GetContext()

	coRpcS, ok := ctx.GetComponent("corpcserver").(*CoRpcServer)
	if ok == true {
		return coRpcS
	}

	host, _ := ctx.CurrentServer["host"].(string)
	port, _ := ctx.CurrentServer["port"].(int)
	rpcS := rpcserver.NewRpcServer(host, port)

	coRpcS = &CoRpcServer{rpcS}
	ctx.RegisteComponent("corpcserver", coRpcS)
	seelog.Info("CoRpcServer create successfully")
	return coRpcS
}
