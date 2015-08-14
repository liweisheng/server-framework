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
	"component/coserver"
	"component/cosession"
	"connector"
	"connector/tcp_connector"
	"context"
	seelog "github.com/cihub/seelog"
	"service/sessionService"
	"strconv"
	"strings"
)

type CoConnector struct {
	ctx    *context.Context
	cnct   connector.Connector
	coserv *coserver.CoServer
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

	cocnct, ok := ctx.GetComponent("coconnector").(*CoConnector)
	if ok == true {
		return cocnct
	}

	var decode func([]byte) (interface{}, error)
	var encode func(string, string, map[string]interface{}) ([]byte, error)

	cnct := getConnector(ctx)

	if opts, ok := ctx.AllOpts["coconnector"]; ok == true {
		decode, _ = opts["decode"].(func([]byte) (interface{}, error))
		encode, _ = opts["encode"].(func(string, string, map[string]interface{}) ([]byte, error))
	}

	coserv, ok1 := ctx.GetComponent("coserver").(*coserver.CoServer)
	if ok1 == false {
		coserv = coserver.NewCoServer()
	}

	cosess, ok2 := ctx.GetComponent("cosession").(*cosession.CoSession)
	if ok2 == false {
		cosess = cosession.NewCoSession()
	}

	coconn, ok3 := ctx.GetComponent("coconnection").(*coconnection.CoConnection)
	if ok3 == false {
		coconn = coconnection.NewCoConnection()
	}

	cocnct = &CoConnector{ctx, cnct, coserv, cosess, coconn, decode, encode}
	ctx.RegisteComponent("coconnector", cocnct)
	return cocnct
}

/// 向sids标示的所有session发送消息.
///
/// @param reqID 请求id
/// @param route 路由
/// @param msg 发送的消息
/// @param sids 接受消息的session
func (cocnct *CoConnector) Send(reqID string, route string, msg map[string]interface{}, sids []uint32) {
	seelog.Debugf("<%v> send msg<%v> with reqID<%v>,route<%v> to sids<%v>", cocnct.ctx.GetServerID(), msg, reqID, route, sids)

	var encodedMsg []byte
	var err error
	if cocnct.encode != nil {
		encodedMsg, err = cocnct.encode(reqID, route, msg)
	} else {
		encodeFunc := cocnct.cnct.Encode
		encodedMsg, err = encodeFunc(reqID, route, msg)
	}

	if err != nil {
		seelog.Errorf("<%v> encode msg<%v> error<%v>", cocnct.ctx.GetServerID(), msg, err.Error())
		return
	}

	for _, sid := range sids {
		go cocnct.cosess.SendMsgBySID(sid, encodedMsg)
	}
}

/// 回调函数，当有新连接到来时调用该回调函数，创建并记录session,该回调函数注册给connector使用.
func (cocnct *CoConnector) ConnectionEventCB(sock connector.Socket) *sessionService.Session {
	session := cocnct.cosess.CreateSession(sock.ID(), cocnct.ctx.GetServerID(), sock)
	return session
}

/// 回调函数，当连接上有新的message到来时调用该回调函数,该回调函数注册给connector使用.
func (cocnct *CoConnector) MessageEventCB(sid uint32, msg map[string]interface{}) {
	route, _ := msg["route"].(string)
	if cocnct.checkRouteValidity(route) == false {
		seelog.Errorf("<%v> invalid route<%v>", cocnct.ctx.GetServerID(), route)
		return
	}
}

/// 检查路由的合法性，路由格式是:serverType.handler.method
///
/// 如考虑请求聊天服务器chatServer新用户登录，则route格式为chatServer.UserManager.
func (cocnct *CoConnector) checkRouteValidity(route string) bool {
	return strings.Index(route, ".") != -1
}

/// 启动组件.
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
	port := strconv.Itoa(curServer["port"].(int))
	return tcp_connector.NewTcpConnector(host, port, nil)
}
