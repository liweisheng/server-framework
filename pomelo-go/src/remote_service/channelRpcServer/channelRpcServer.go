package channelRpcServer

import (
	"component/coconnector"
	"component/corpcserver"
	"component/cosession"
	"context"
	"errors"
	seelog "github.com/cihub/seelog"
	"os"
)

type ChannelRpcServer struct {
	coCnct *coconnector.CoConnector
}

func (crs *ChannelRpcServer) Start() {
	ctx := context.GetContext()
	coRpcS, ok := ctx.GetComponent("corpcserver").(*corpcserver.CoRpcServer)

	if ok == false {
		coRpcS = corpcserver.NewCoRpcServer()
	}

	seelog.Info("ChannelRpcServer start,registe service to rpcserver")
	// srs := NewSessionRpcServer()
	coRpcS.RegisteService(crs)

}

func NewChannelRpcServer() *ChannelRpcServer {
	ctx := context.GetContext()

	coCnct, ok := ctx.GetComponent("coconnector").(*coconnector.CoConnector)
	if ok == false {
		coCnct = coconnector.NewCoConnector(ctx)
	}

	return &ChannelRpcServer{coCnct}
}

func (crs *ChannelRpcServer) PushMessage(args map[string]interface{}, reply interface{}) error {
	route, _ := args["route"].(string)
	msg, ok1 := args["msg"].(map[string]interface{})
	if ok1 == false {
		seelog.Error("Invalid or empty message")
		return errors.New("Invalid or empty message")
	}

	coSess, ok2 := context.GetContext().GetComponent("cosession").(*cosession.CoSession)

	if ok2 == false {
		seelog.Critical("Failed to get component<cosession>")
		os.Exit(1)
	}

	uids, ok3 := args["uids"].([]string)
	if ok3 == false {
		seelog.Error("Failed to get uids")
		return errors.New("Failed to get uids")
	}

	sids := make([]uint32, 0, 16)
	for _, uid := range uids {
		sessions := coSess.GetSessionsByUID(uid)

		if sessions == nil {
			seelog.Warnf("Failed to push message to uid<%v> because of nil sessions", uid)
			continue
		}

		for _, session := range sessions {
			sids = append(sids, session.Id)
		}
	}

	seelog.Debugf("<%v> push messages<%v> to uids<%v> with sessions<%v>", context.GetContext().GetServerID(), msg, uids, sids)

	//XXX:调用coconnector的send方法发送.
	crs.coCnct.Send("", route, msg, sids)
	return nil

}

func (crs *ChannelRpcServer) Broadcast(args map[string]interface{}, reply interface{}) error {
	//XXX:调用coconnector的send方法发送.
	return nil
}
