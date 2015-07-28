/**
 * author:liweisheng date:2015/07/26
 */

/**
 * backendsession 相当于前端服务器session在后端的代理，仅由后端服务器加载.
 * backendsession 通过rpc与前端session交互.
 */

package backendSessionService

import (
	"component/corpcclient"
	"context"
	seelog "github.com/cihub/seelog"
)

type BackendSessionService struct {
	rpcCient *corpcclient.CoRpcClient
}

func NewBackendSessionService(ctx *context.Context) *BackendSessionService {

	rpcClient, ok := ctx.GetComponent("corpcclient").(*corpcclient.CoRpcClient)

	if ok == false || rpcClient == nil {
		rpcClient = corpcclient.NewCoRpcClient()
	}

	return &BackendSessionService{rpcClient}
}

/// 创建新的BackendSession.
func (bss *BackendSessionService) CreateBackendSession(opts map[string]interface{}) *BackendSession {
	return newBackendSession(opts, bss)
}

/// 通过frontendID 和sessionID获得前端服务器session的信息.
///
/// @param frontendID 前端服务器id
/// @param sid session id
/// @return backend session
func (bss *BackendSessionService) GetBackendSessionBySID(frontendid string, sid uint32) *BackendSession {
	seelog.Tracef("frontendid<%v>,sid<%v>", frontendid, sid)

	reply := make(map[string]interface{})
	method := "SessionRpcServer.GetSessionBySID"
	rpcRelpy := bss.rpcCient.RpcCall(frontendid, method, sid, &reply)

	if rpcRelpy != nil && rpcRelpy.Error == nil {
		<-rpcRelpy.Done

		if reply == nil {
			return nil
		}
		opts, _ := reply["opts"].(map[string]interface{})

		backendSession := bss.CreateBackendSession(opts)
		backendSession.uid = reply["uid"].(string)
		backendSession.id = reply["sid"].(uint32)
		backendSession.frontendID = reply["frontendid"].(string)

		seelog.Debugf("Receive from rpc server<%v>", reply)
		return backendSession
	} else {
		seelog.Error("Rpc Call failed")
		return nil
	}
}

/// 通过user id 获得所有绑定的session
///
/// @param frontendid 前端服务器id
/// @param uid user id
func (bss *BackendSessionService) GetBackendSessionsByUID(frontendid string, uid string) []*BackendSession {
	seelog.Tracef("frontendid<%v>,sid<%v>", frontendid, uid)

	backendSessions := make([]*BackendSession, 0, 5)

	replies := make([]map[string]interface{}, 0, 5)
	method := "SessionRpcServer.GetSessionsByUID"

	rpcReply := bss.rpcCient.RpcCall(frontendid, method, uid, &replies)

	if rpcReply != nil && rpcReply.Error == nil {
		<-rpcReply.Done

		if replies == nil {
			return nil
		}
		for _, elem := range replies {

			opts, _ := elem["opts"].(map[string]interface{})
			seelog.Debugf("Replied session opts<%v>", opts)

			bs := bss.CreateBackendSession(opts)

			bs.uid, _ = elem["uid"].(string)
			bs.id, _ = elem["sid"].(uint32)
			bs.frontendID, _ = elem["frontendid"].(string)
			backendSessions = append(backendSessions, bs)
		}

		return backendSessions
	} else {
		return nil
	}

}

/// 通过session id从前端服务器踢除用户连接，通过rpc调用前端服务器的同样操作.
///
/// @param frontendid 前端服务器id
/// @param sid session id
/// @reason 作为踢除用户连接时发送给用户端的提示信息.
/// @return 无返回值
func (bss *BackendSessionService) KickBySID(frontendid string, sid uint32, reason string) {
	seelog.Tracef("frontendid<%v>,sid<%v>", frontendid, sid)

	method := "SessionRpcServer.KickBySID"

	args := make([]interface{}, 2)

	args[0] = sid
	args[1] = reason
	if reply := bss.rpcCient.RpcCall(frontendid, method, args, nil); reply != nil {

		<-reply.Done

		if reply.Error != nil {
			seelog.Error(reply.Error.Error())
		}
	}
}

/// 通过uid从前端服务器剔除用户，通过rpc调用前端服务器的同样操作.
///
/// @param frontendid 前端服务器id
/// @param uid 用户id
/// @param reason 作为踢出用户前发送的提示信息
func (bss *BackendSessionService) KickByUID(frontendid string, uid string, reason string) {
	seelog.Tracef("frontendif<%v>,uid<%v>", frontendid, uid)

	method := "SessionRpcServer.KickByUID"

	args := make([]string, 2)

	args[0] = uid
	args[1] = reason

	if reply := bss.rpcCient.RpcCall(frontendid, method, args, nil); reply != nil {
		<-reply.Done

		if reply.Error != nil {
			seelog.Error(reply.Error.Error())
		}
	}
}

/// 绑定用户id与session id，该操作会通过rpc调用在前端服务器完成uid与sid的绑定.
///
/// @param frontendid 前端服务器id
/// @param sid session id
/// @param uid 用户id
func (bss *BackendSessionService) BindUID(frontendid string, sid uint32, uid string) {
	seelog.Tracef("Bind uid<%v> with session id<%v>", uid, sid)

	args := make([]interface{}, 2)

	args[0] = sid
	args[1] = uid

	method := "SessionRpcServer.BindUID"

	if reply := bss.rpcCient.RpcCall(frontendid, method, args, nil); reply != nil {
		<-reply.Done

		if reply.Error != nil {
			seelog.Error(reply.Error.Error())
		}
	}
}

/// 解除用户id到session id的绑定，该操作会通过rpc调用影响前端服务器.
///
/// @param frontendid 前端服务器id
/// @param sid session id
/// @param uid 用户id
func (bss *BackendSessionService) UnbindUID(frontendid string, sid uint32, uid string) {
	seelog.Tracef("Unbind uid<%v> with session id<%v>", uid, sid)

	args := make([]interface{}, 2)

	args[0] = sid
	args[1] = uid

	method := "SessionRpcServer.UnbindUID"

	if reply := bss.rpcCient.RpcCall(frontendid, method, args, nil); reply != nil {
		<-reply.Done

		if reply.Error != nil {
			seelog.Error(reply.Error.Error())
		}
	}
}

/// 将设置sid标识的session的属性同步到前端服务器.
///
/// @param frontendid 前端服务器id
/// @param sid session id
/// @param key 属性名
/// @param value 属性值
func (bss *BackendSessionService) PushOpt(frontendid string, sid uint32, key string, value interface{}) {
	seelog.Tracef("Push opt to <%v> with sid<%v>,key<%v>,value<%v>", frontendid, sid, key, value)

	args := make(map[string]interface{})

	args["sid"] = sid
	args["key"] = key
	args["value"] = value

	method := "SessionRpcServer.PushOpt"

	if reply := bss.rpcCient.RpcCall(frontendid, method, args, nil); reply != nil {
		<-reply.Done

		if reply.Error != nil {
			seelog.Error(reply.Error.Error())
		}
	}

}

func (bss *BackendSessionService) PushAllOpts(frontendid string, sid uint32, settings map[string]interface{}) {
	seelog.Tracef("Push all opts to <%v> with sid<%v>", frontendid, sid)

	args := make(map[string]interface{})

	args["sid"] = sid
	args["settings"] = settings

	method := "SessionRpcServer.PushAllOpts"

	if reply := bss.rpcCient.RpcCall(frontendid, method, args, nil); reply != nil {
		<-reply.Done

		if reply.Error != nil {
			seelog.Error(reply.Error.Error())
		}
	}
}

type BackendSession struct {
	frontendID            string
	uid                   string
	id                    uint32
	backendSessionService *BackendSessionService
	opts                  map[string]interface{}
}

/// 创建新的BackendSession,外部包不能访问此方法,创建BackendSession参见CreateBackendSession.
///
/// @param opts BackendSession选项
/// @param backendSessionService 管理BackendSession的BackendSessionService,及创建BackendSession者.
func newBackendSession(opts map[string]interface{}, backendSessionService *BackendSessionService) *BackendSession {
	return &BackendSession{"", "", 0, backendSessionService, opts}
}

func (bs *BackendSession) BindUID(uid string) {
	bs.backendSessionService.BindUID(bs.frontendID, bs.id, uid)
}

func (bs *BackendSession) UnbindUID(uid string) {
	bs.backendSessionService.UnbindUID(bs.frontendID, bs.id, uid)
}

func (bs *BackendSession) SetOpt(key string, value interface{}) {
	bs.opts[key] = value
}

func (bs *BackendSession) GetOpt(key string) interface{} {
	return bs.opts[key]
}

func (bs *BackendSession) PushOpt(key string) {
	bs.backendSessionService.PushOpt(bs.frontendID, bs.id, key, bs.opts[key])
}

func (bs *BackendSession) PushAllOpts() {
	bs.backendSessionService.PushAllOpts(bs.frontendID, bs.id, bs.opts)
}
