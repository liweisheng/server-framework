package coconnection

import (
	"context"
	"github.com/cihub/seelog"
	"service/connectionService"
)

type CoConnection struct {
	*connectionService.ConnectionService
	ctx *context.Context
}

func NewCoConnection() *CoConnection {
	ctx := context.GetContext()
	coconn, ok := ctx.GetComponent("coconnection").(*CoConnection)
	if ok == true {
		return coconn
	}

	cs := connectionService.NewConnectionService(ctx.GetServerID())

	coconn = &CoConnection{cs, ctx}
	ctx.RegisteComponent("coconnection", coconn)

	seelog.Infof("CoConnetion create successfully")
	return coconn
}
