package coconnection

import (
	"context"
	"service/connectionService"
)

type CoConnection struct {
	*connectionService.ConnectionService
	ctx *context.Context
}

func NewCoConnection(ctx *context.Context) *CoConnection {
	cs := connectionService.NewConnectionService(ctx.GetServerID())

	coconn := &CoConnection{cs, ctx}
	ctx.RegisteComponent("coconnection", coconn)
	return coconn
}
