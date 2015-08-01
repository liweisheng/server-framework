/**
 * author:liweisheng  date:2015/07/01
 */
package server

import (
	"component/corpcclient"
	"component/corpcserver"
	"context"
	"github.com/cihub/seelog"
	"reflect"
	"service/sessionService"
	"strings"
)

type Server struct {
	ctx      *context.Context
	coRpcS   *corpcserver.CoRpcServer
	coRpcC   *corpcclient.CoRpcClient
	handlers map[string]func(*sessionService.SimpleSession, map[string]interface{}) map[string]interface{}
}

func NewServer() *Server {
	ctx := context.GetContext()

	coRpcS, ok1 := ctx.GetComponent("corpcserver").(*corpcserver.CoRpcServer)
	if ok1 == false {
		coRpcS = corpcserver.NewCoRpcServer()
	}

	coRpcC, ok2 := ctx.GetComponent("corpcclient").(*corpcclient.CoRpcClient)
	if ok2 == false {
		coRpcC = corpcclient.NewCoRpcClient()
	}
	handlers := make(map[string]func(*sessionService.SimpleSession, map[string]interface{}) map[string]interface{})
	return &Server{ctx, coRpcS, coRpcC, handlers}
}

func (s *Server) RegisteAsLocalService(recv interface{}) {
	t := reflect.TypeOf(recv)
	v := reflect.ValueOf(recv)
	numMethod := t.NumMethod()
	sName := v.Type().Name()
	for index := 0; index < numMethod; index++ {
		methodName := t.Method(index).Name
		if handler, ok := v.Method(index).Interface().(func(*sessionService.SimpleSession, map[string]interface{}) map[string]interface{}); ok == true {
			s.handlers[sName+"."+methodName] = handler
			seelog.Infof("<%v> Registe local service<%v>'s method<%v>", context.GetContext().GetServerID(), v.Type().Name(), methodName)
		} else {
			seelog.Warnf("<%v> Registe local service <%v>'s method<%v> error", context.GetContext().GetServerID(), v.Type().Name(), methodName)
		}
	}
}

func (s *Server) RegisteAsRemoteService(recv interface{}) {
	s.coRpcS.RegisteService(recv)
}

func (s *Server) GlobalHandler(msg map[string]interface{}, session *sessionService.SimpleSession, sendCB func(string, string, map[string]interface{}, []uint32)) {
	frontendID := s.ctx.GetServerID()
	route, _ := msg["route"].(string)
	record := parseRoute(route)
	if record == nil {
		seelog.Errorf("receive invalid route<%v>, route format must be<serverType.handler.method>", route)
		return
	}
	handlerMethodName := record.handler + "." + record.method
	if context.GetContext().GetServerType() == record.serverType {
		///调用本地方法处理.
		seelog.Infof("<%v> receive local handler<%v> request", frontendID, handlerMethodName)
		if handlerMethod, ok := s.handlers[handlerMethodName]; ok == true {
			replyMsg := handlerMethod(session, msg)
			seelog.Debugf("<%v> local handler<%v> reply msg<%v>", frontendID, handlerMethodName, replyMsg)
			sids := make([]uint32, 1)
			sids[0] = session.Id
			reqID, _ := replyMsg["reqID"].(string)
			rroute, _ := replyMsg["route"].(string)
			msgBody, _ := replyMsg["msg"].(map[string]interface{})
			sendCB(reqID, rroute, msgBody, sids)
		} else {
			seelog.Errorf("<%v> request local handler<%v> error", frontendID, handlerMethodName)
		}
	} else {
		///rpc调用远端方法

		servID := s.ctx.GetServerIDByType(record.serverType)
		seelog.Infof("<%v> invoke remote handler<%v> request of <%v>", frontendID, handlerMethodName, servID)
		reply := make(map[string]interface{})
		if err := s.coRpcC.RpcCall(servID, handlerMethodName, msg, &reply); err != nil {
			seelog.Errorf("<%v> rpc call<%v> error<%v>", frontendID, handlerMethodName, err.Error())
		}
	}

}

func (s *Server) Start() {
	seelog.Infof("<%v> component server start...", s.ctx.GetServerID())
}

type routeRecord struct {
	route      string
	serverType string
	handler    string
	method     string
}

func parseRoute(route string) *routeRecord {
	rs := strings.Split(route, ".")
	if len(rs) != 3 {
		return nil
	}

	return &routeRecord{route, rs[0], rs[1], rs[2]}

}
