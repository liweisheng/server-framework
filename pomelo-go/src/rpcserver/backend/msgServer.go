/**
 * author:liweisheng  date:2015/07/23
 */

/**
 * msgServer是由后端服务器加载的rpc服务器端，它接收并服务前端服务器发起的rpc调用请求.
 * 采用json做编码/解码协议.
 */

package msgServer

import (
	seelog "github.com/cihub/seelog"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

type MsgServer struct {
	host string
	port string
	rcvr interface{}
}

func NewMsgServer(host string, port string, rcvr interface{}) *MsgServer {
	return &MsgServer{host, port, rcvr}
}

func (ms *MsgServer) RegisteService(r interface{}) {
	ms.rcvr = r
}

func (ms *MsgServer) Start() {

}
