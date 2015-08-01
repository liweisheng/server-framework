package coserver

import (
	"context"
	"github.com/cihub/seelog"
	"server"
)

type CoServer struct {
	*server.Server
}

func NewCoServer() *CoServer {
	ctx := context.GetContext()
	coServ, ok := ctx.GetComponent("coserver").(*CoServer)
	if ok == true {
		return coServ
	}
	serv := server.NewServer()
	coServ = &CoServer{serv}
	seelog.Infof("<%v> component server created", ctx.GetServerID())
	return coServ
}
