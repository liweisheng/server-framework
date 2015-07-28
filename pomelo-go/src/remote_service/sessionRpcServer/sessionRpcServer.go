package sessionRpcServer

import (
	"component/corpcserver"
	"component/cosession"
	"context"
	seelog "github.com/cihub/seelog"
)

type SessionRpcServer struct {
	sessionService *cosession.CoSession
}

func Start() {
	ctx := context.GetContext()
	coRpcS, ok := ctx.GetComponent("corpcserver").(corpcserver.CoRpcServer)

	if ok == true {
		seelog.Info("SessionRpcServer start")
		srs := NewSessionRpcServer()
		coRpcS.RegisteService(srs)
	} else {
		seelog.Error("SessionRpcServer failed to start")
	}
}

func NewSessionRpcServer() *SessionRpcServer {
	ctx := context.GetContext()

	sessionService, ok := ctx.GetComponent("cosession").(*cosession.CoSession)

	if ok == false || sessionService == nil {
		sessionService = cosession.NewCoSession(ctx)
	}

	return &SessionRpcServer{sessionService}
}

func (srs *SessionRpcServer) GetSessionBySID(sid uint32, reply *map[string]interface{}) error {

	seelog.Debugf("SessionRpcServer method<GetSessionBySID> is invoked with sid<%v>", sid)
	session := srs.sessionService.GetSessionByID(sid)

	if session != nil {
		(*reply)["sid"] = sid
		(*reply)["uid"] = session.Uid
		(*reply)["frontendid"] = session.FrontendID
		(*reply)["opts"] = session.Opts
	} else {
		*reply = nil
	}
	return nil
}

func (srs *SessionRpcServer) GetSessionsByUID(uid string, reply *[]map[string]interface{}) error {
	seelog.Debugf("SessionRpcServer method<GetSessionByUID> is invoked with uid<%v>", uid)

	sessions := srs.sessionService.GetSessionsByUID(uid)

	if sessions != nil {
		for _, elem := range sessions {

			s := make(map[string]interface{})
			s["sid"] = elem.Id
			s["uid"] = uid
			s["frontendid"] = elem.FrontendID
			s["opts"] = elem.Opts
			*reply = append(*reply, s)
		}
	} else {
		*reply = nil
	}

	return nil
}

func (srs *SessionRpcServer) KickBySID(args []interface{}, reply interface{}) error {

	sid, _ := args[0].(uint32)
	reason, _ := args[1].(string)

	seelog.Debugf("SessionRpcServer method<KickBySID> is invoked with sid<%v> reason<%v>", sid, reason)
	srs.sessionService.KickBySessionID(sid, reason)

	return nil
}

func (srs *SessionRpcServer) KickByUID(args []string, reply interface{}) error {

	uid := args[0]
	reason := args[1]

	seelog.Debugf("SessionRpcServer method<KickByUID> is invoked with uid<%v> reason<%v>", uid, reason)

	srs.sessionService.KickByUID(uid, reason)

	return nil
}

func (srs *SessionRpcServer) BindUID(args []interface{}, reply interface{}) error {

	sid, _ := args[0].(uint32)
	uid, _ := args[1].(string)

	seelog.Debugf("SessionRpcServer method<BindUID> is invoked with sid<%v> uid<%v>", sid, uid)

	srs.sessionService.BindUID(uid, sid)

	return nil
}

func (srs *SessionRpcServer) UnbindUID(args []interface{}, reply interface{}) error {
	sid, _ := args[0].(uint32)
	uid, _ := args[1].(string)

	seelog.Debugf("SessionRpcServer method<UnbindUID> is invoked with sid<%v> uid<%v>", sid, uid)

	srs.sessionService.UnbindUID(uid, sid)

	return nil
}

func (srs *SessionRpcServer) PushOpt(args map[string]interface{}, reply interface{}) error {

	sid, _ := args["sid"].(uint32)
	key, _ := args["key"].(string)
	value, _ := args["value"]

	seelog.Debugf("SessionRpcServer method<PushOpt> is invoked with sid<%v> key<%v> value<%v>", sid, key, value)

	srs.sessionService.PushOpt(sid, key, value)

	return nil
}

func (srs *SessionRpcServer) PushAllOpts(args map[string]interface{}, reply interface{}) error {

	sid, _ := args["sid"].(uint32)
	opts, _ := args["settings"].(map[string]interface{})

	seelog.Debugf("SessionRpcServer method<PushAllOpt> is invoked with sid<%v>", sid)

	srs.sessionService.PushAllOpts(sid, opts)

	return nil
}
