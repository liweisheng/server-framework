package corpcclient

import (
	"context"
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

	return coRpcC

}
