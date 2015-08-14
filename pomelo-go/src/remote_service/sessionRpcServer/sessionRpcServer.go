/**
 * author:liweisheng date:2015/07/24
 */

/**
 * sessionRpcServer包与包backendSessionService相对等，相当于rpc服务的服务端与客户端，backendSessionService
 * 发起rpc调用请求sessionRpcService服务，将对backendsession的改变影响到前端服务器的session.
 *
 * sessionRpcServer内部想象session的操作都是用过组建cosession实现的.
 */
package sessionRpcServer

import (
	"component/corpcserver"
	"component/cosession"
	"context"
	"fmt"

	seelog "github.com/cihub/seelog"
)

type NoReply int
type SessionRpcServer struct {
	sessionService *cosession.CoSession
}

func (srs *SessionRpcServer) Start() {
	ctx := context.GetContext()
	coRpcS, ok := ctx.GetComponent("corpcserver").(*corpcserver.CoRpcServer)

	if ok == false {
		coRpcS = corpcserver.NewCoRpcServer()
	}

	seelog.Infof("frontendserver<%v> SessionRpcServer start,registe service to rpc server", ctx.GetServerID())
	// srs := NewSessionRpcServer()
	coRpcS.RegisteService(srs)
}

func NewSessionRpcServer() *SessionRpcServer {
	ctx := context.GetContext()

	sessionService, ok := ctx.GetComponent("cosession").(*cosession.CoSession)

	if ok == false || sessionService == nil {
		sessionService = cosession.NewCoSession()
	}

	return &SessionRpcServer{sessionService}
}

/// rpc服务端操作，通过session id返回前端session的详细信息.
///
/// @param sid session id
/// @param reply 返回给rpc调用端，name->value形式存放,存放信息包括,sid:session id, uid:用户id,frontendid:前端服务器id,opts:设置的属性
/// @return nil
func (srs *SessionRpcServer) GetSessionBySID(sid *uint32, reply *map[string]interface{}) error {

	seelog.Debugf("SessionRpcServer method<GetSessionBySID> is invoked with sid<%v>", sid)
	session := srs.sessionService.GetSessionByID(*sid)

	rep := make(map[string]interface{})

	if session != nil {
		rep["sid"] = uint32(*sid)
		rep["uid"] = session.Uid
		rep["frontendid"] = session.FrontendID
		rep["opts"] = session.Opts
		*reply = rep
	} else {
		*reply = rep
	}
	return nil
}

/// 通过用户id返回用户id绑定的所有session(如果容许同一个用户id绑定多次，则存在多个session)
///
/// @param uid 用户id
/// @param reply 返回给rpc 客户端的结果，有多个session组成的数组，每个数组元素以name->value形式存放.
func (srs *SessionRpcServer) GetSessionsByUID(uid string, reply *[]map[string]interface{}) error {
	seelog.Debugf("SessionRpcServer method<GetSessionByUID> is invoked with uid<%v>", uid)

	sessions := srs.sessionService.GetSessionsByUID(uid)
	rep := make([]map[string]interface{}, 0)
	if sessions != nil {
		for _, elem := range sessions {

			s := make(map[string]interface{})
			s["sid"] = elem.Id
			s["uid"] = uid
			s["frontendid"] = elem.FrontendID
			s["opts"] = elem.Opts
			rep = append(rep, s)
		}
		*reply = rep
	} else {
		*reply = nil
	}
	return nil
}

/// rpc服务端操作,通过制定session id踢除用户.
///
/// @param args args[0]表示sessionid{uint32} ，args[1]表示reason{string}
/// @param reply 无实际意义参数，rpc客户端传来的nil
/// @return nil
func (srs *SessionRpcServer) KickBySID(args []interface{}, reply *NoReply) error {

	sid, _ := args[0].(float64)
	reason, _ := args[1].(string)

	seelog.Debugf("SessionRpcServer method<KickBySID> is invoked with sid<%v> reason<%v>", sid, reason)
	srs.sessionService.KickBySessionID(uint32(sid), reason)

	return nil
}

/// rpc服务端操作，通过用户id踢出用户.
/// @param args args[0]{string}用户id, args[1]{string} reason
/// @param reply nil
/// @return nil
func (srs *SessionRpcServer) KickByUID(args []string, reply *NoReply) error {

	uid := args[0]
	reason := args[1]

	seelog.Debugf("SessionRpcServer method<KickByUID> is invoked with uid<%v> reason<%v>", uid, reason)

	srs.sessionService.KickByUID(uid, reason)

	return nil
}

/// rpc服务端操作, 将用户id与session id绑定.
///
/// @param args arg[0]{uint32} sessionid, args[1]{string}用户id
/// @param reply 当前rpc客户端传过来的nil
/// @return nil
func (srs *SessionRpcServer) BindUID(args *[]interface{}, reply *NoReply) error {

	sid, _ := (*args)[0].(float64)
	uid, _ := (*args)[1].(string)
	fmt.Printf("BindUID args:%v,sid:%v \n", *args, uint32(sid))
	seelog.Debugf("SessionRpcServer method<BindUID> is invoked with sid<%v>, uid<%v>", sid, uid)

	srs.sessionService.BindUID(uid, uint32(sid))

	return nil
}

/// rpc服务端操作， 将用户id与session id解绑定.
///
/// @param args args[0]{uint32} sessionid, args[1] {string} 用户id
/// @param reply 当前rpc客户端传过来为nil
/// @return nil
func (srs *SessionRpcServer) UnbindUID(args *[]interface{}, reply *NoReply) error {
	sid, _ := (*args)[0].(float64)
	uid, _ := (*args)[1].(string)

	seelog.Debugf("SessionRpcServer method<UnbindUID> is invoked with sid<%v> uid<%v>", sid, uid)

	srs.sessionService.UnbindUID(uid, uint32(sid))

	return nil
}

/// rpc服务端操作，设置用户的属性.
///
/// @param args args["sid"]{uint32}session id， args["key"]属性名，args["value"]{interface{}}属性值
/// @param reply rpc客户端传递多来为nil
/// @return nil
func (srs *SessionRpcServer) PushOpt(args *map[string]interface{}, reply *NoReply) error {

	sid, _ := (*args)["sid"].(float64)
	key, _ := (*args)["key"].(string)
	value, _ := (*args)["value"]

	seelog.Debugf("SessionRpcServer method<PushOpt> is invoked with sid<%v> key<%v> value<%v>", sid, key, value)

	srs.sessionService.PushOpt(uint32(sid), key, value)

	return nil
}

/// rpc服务端操作,类似PushOpt, 可以同时设置多个属性.
///
/// @param args["sid"]{uint32} sessionid，args["setting"] {map[string]interface{}}多个属性的key->value映射.
/// @param reply rpc客户端传递过来为nil
/// @return nil
func (srs *SessionRpcServer) PushAllOpts(args *map[string]interface{}, reply *NoReply) error {

	sid, _ := (*args)["sid"].(float64)
	opts, _ := (*args)["settings"].(map[string]interface{})

	seelog.Debugf("SessionRpcServer method<PushAllOpt> is invoked with sid<%v>", sid)

	srs.sessionService.PushAllOpts(uint32(sid), opts)

	return nil
}
