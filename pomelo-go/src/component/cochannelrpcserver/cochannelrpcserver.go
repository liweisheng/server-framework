package cochannelrpcserver

import (
	"context"
	"github.com/cihub/seelog"
	"remote_service/channelRpcServer"
)

type CoChannelRpcServer struct {
	*channelRpcServer.ChannelRpcServer
}

func NewCoChannelRpcServer() *CoChannelRpcServer {

	ctx := context.GetContext()

	coChanRS, ok := ctx.GetComponent("cochannelrpcserver").(*CoChannelRpcServer)
	if ok == false {
		seelog.Infof("CoChannelRpcServer not found,create new...")
		chanRpcS := channelRpcServer.NewChannelRpcServer()
		coChanRS = &CoChannelRpcServer{chanRpcS}
		ctx.RegisteComponent("cochannelrpcserver", coChanRS)
	}

	seelog.Infof("CoChannelRpcServer create successfully")
	return coChanRS
}
