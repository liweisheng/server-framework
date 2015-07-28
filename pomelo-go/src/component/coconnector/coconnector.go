/**
 * author:liweisheng date:2015/07/14
 */

/**
 * CoConnector组件由前端服务器加载，负责接收客户端连接，维护所有到来的连接信息，接收客户端请求
 * 并响应客户端请求。并对不同的用户请求做出不同处理：
 *    （1）如果用户请求的是前端服务，则前端服务器
 * 	  （2）如果用户请求的不是

 */

package coconnector

import (
	"component/coconnection"
	"component/cosession"
	"connector"
	"connector/tcp_connector"
	"context"
)

type CoConnector struct {
	ctx    *context.Context
	cnct   connector.Connector
	cosess *cosession.CoSession
	coconn *coconnection.CoConnection
	decode func([]byte) (interface{}, error)
	encode func(string, string, map[string]interface{}) ([]byte, error)
}

/// 创建新的CoConnetor组件.
///
/// 创建CoConnector组件时，CoConnector组件使用到的CoSession组件,CoConnection组件一同创建，
/// CoServer组件则通过Context拿到，在CoConnector组件启动时加载.
func NewCoConnector(ctx *context.Context) *CoConnector {
	var decode func([]byte) (interface{}, error)
	var encode func(string, string, map[string]interface{}) ([]byte, error)

	cnct := getConnector(ctx)

	if opts, ok := ctx.AllOpts["coconnector"]; ok == true {
		decode, _ = opts["decode"].(func([]byte) (interface{}, error))
		encode, _ = opts["encode"].(func(string, string, map[string]interface{}) ([]byte, error))
	}

	cosess, ok := ctx.GetComponent("cosession").(*cosession.CoSession)

	coconn, _ := ctx.GetComponent("coconnection").(*coconnection.CoConnection)

	cocnct := &CoConnector{ctx, cnct, cosess, coconn, decode, encode}
	ctx.RegisteComponent("coconnector", cocnct)
	return cocnct
}

func (cc *CoConnector) Start() {

}

/// 获取CoConnector组件内部用于发送和接收消息的connector服务.
///
/// 如果用户配置了connector(保存在Context中)，则使用全局配置的connector,
/// 否者使用默认的TcpConnector.
func getConnector(ctx *context.Context) connector.Connector {
	opts := ctx.AllOpts["coconnector"]
	if opts != nil {
		if c, setted := opts["connector"]; setted == true {
			if cnct, ok := c.(connector.Connector); ok == true {
				return cnct
			}
		}
	}

	curServer := ctx.CurrentServer
	host, _ := curServer["host"].(string)
	port, _ := curServer["port"].(string)
	return tcp_connector.NewTcpConnector(host, port, nil)
}
