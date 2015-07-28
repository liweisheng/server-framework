package cosessionrpcserver

import (
	"context"
	"remote_service/sessionRpcServer"
)

type CoSessionRpcServer struct {
	*sessionRpcServer.SessionRpcServer
}

func NewCoSessionRpcServer() *CoSessionRpcServer {
	ctx := context.GetContext()
	coSRS, ok := ctx.GetComponent("cosessionrpcserver").(*CoSessionRpcServer)

	if ok == true {
		return coSRS
	}

	srs := sessionRpcServer.NewSessionRpcServer()
	coSRS = &CoSessionRpcServer{srs}

	ctx.RegisteComponent("cosessionrpcserver", coSRS)

	return coSRS
}
