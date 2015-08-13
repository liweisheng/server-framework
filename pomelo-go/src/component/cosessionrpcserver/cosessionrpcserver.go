/**
 * author: liweisheng date:2015/07/25
 */

/// 组建cosessionrpcserver 内部封装sessionRpcServer, 详细介绍参考remote_service/sessionRpcServer.go
package cosessionrpcserver

import (
	"context"
	"remote_service/sessionRpcServer"
)

type CoSessionRpcServer struct {
	*sessionRpcServer.SessionRpcServer
}

/// 创建CoSessionRpcServer组建.
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
