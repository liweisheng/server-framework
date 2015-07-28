package corpcserver

import (
	"context"
	"rpcserver"
)

type CoRpcServer struct {
	*rpcserver.RpcServer
}

func NewCoRpcServer() *CoRpcServer {
	ctx := context.GetContext()

	host, _ := ctx.CurrentServer["host"].(string)
	port, _ := ctx.CurrentServer["port"].(int)
	rpcS := rpcserver.NewRpcServer(host, port)

	coRpcS := &CoRpcServer{rpcS}
	ctx.RegisteComponent("corpcserver", coRpcS)

	return coRpcS
}
