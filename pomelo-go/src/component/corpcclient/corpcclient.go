package corpcclient

import (
	"context"
	"github.com/cihub/seelog"
	"rpcclient"
)

type CoRpcClient struct {
	*rpcclient.RpcClient
}

func NewCoRpcClient() *CoRpcClient {
	ctx := context.GetContext()

	coRpcC, ok := ctx.GetComponent("corpcclient").(*CoRpcClient)

	if ok == true {
		return coRpcC
	}

	rpcC := rpcclient.NewRpcClient(ctx)
	coRpcC = &CoRpcClient{rpcC}

	ctx.RegisteComponent("corpcclient", coRpcC)
	seelog.Info("CoRpcClient create successfully")
	return coRpcC

}
