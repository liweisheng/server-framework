/**
 * author:liweisheng date:2015/07/20
 */

/**
 * rpcClient实现rpc客户端,负责
 */

package rpcclient

import (
	"context"
	"fmt"
	seelog "github.com/cihub/seelog"

	// "net/rpc"
	"net/rpc/jsonrpc"
)

type RpcClient struct {
	ctx *context.Context
}

func NewRpcClient(ctx *context.Context) *RpcClient {
	return &RpcClient{ctx}
}

/// 发起远程调用.
///
/// @param serverID 接收远程调用的服务器id
/// @param method 请求的方法
/// @param args 传递给远程过程的参数
/// @param reply 用于传输返回值,为空时表示不关心返回值.
/// @return 当参数reply不为空时，返回值意义参见rpc.Go
func (rc *RpcClient) RpcCall(serverID string, method string, args interface{}, reply interface{}) error {
	info := rc.ctx.GetServerInfoByID(serverID)

	host, _ := info["host"].(string)
	port, _ := info["port"].(int)

	hostPort := fmt.Sprintf("%v:%v", host, port)

	client, err := jsonrpc.Dial("tcp", hostPort)

	if err != nil {
		seelog.Errorf("Fail to Dial rpc server,error message:%v", err.Error())
		return nil
	}

	/// BUG：以下全是BUG
	defer client.Close()
	seelog.Debugf("Rpc client call for remote method<%v>,remote serverid<%v>", method, serverID)
	return client.Call(method, args, reply)

}
