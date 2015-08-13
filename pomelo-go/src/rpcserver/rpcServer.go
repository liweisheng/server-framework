/**
 * author:liweisheng  date:2015/07/23
 */

/**
 * rpcServer是由后端服务器加载的rpc服务器端，它接收并服务前端服务器发起的rpc调用请求.
 * 采用json做编码/解码协议.
 */

package rpcserver

import (
	"fmt"
	seelog "github.com/cihub/seelog"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

type RpcServer struct {
	host      string
	port      int
	rpcServer *rpc.Server
}

func NewRpcServer(host string, port int) *RpcServer {
	rs := rpc.NewServer()
	return &RpcServer{host, port, rs}
}

/// 向rpc服务器注册服务,封装rpc.Register
func (ms *RpcServer) RegisteService(r interface{}) {
	err := ms.rpcServer.Register(r)
	if err != nil {
		seelog.Criticalf("Fail to Register Rpc Service,%v", err.Error())
		os.Exit(1)
	}
}

/// 启动rpcServer,监听rpc服务器端口,由于Start内部调用阻塞的方法,应在go 语句中调用.
func (ms *RpcServer) Start() {
	go func() {
		seelog.Info("RpcServer start...")
		hostAndPort := fmt.Sprintf("%v:%v", ms.host, ms.port)

		servAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)

		if err != nil {
			seelog.Criticalf("RpcServer failed to start with err<%v>", err.Error())
			os.Exit(1)
		}

		listener, err := net.ListenTCP("tcp4", servAddr)

		if err != nil {
			seelog.Criticalf("RpcServer failed to start with err<%v>", err.Error())
			os.Exit(1)
		}

		seelog.Debugf("Rpc Server listening: <%v>", servAddr.String())
		defer listener.Close()

		for {
			conn, err := listener.Accept()

			seelog.Debug("Rpc Server accept new connection")
			if err != nil {
				seelog.Critical(err.Error())
				os.Exit(1)
			}
			go ms.rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}()

}
