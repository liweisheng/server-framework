/// 测试backendSessionService,backendSessionService中封装了对backendsession的操作，对backendsession的操作
/// 会通过发起对sessionRpcServer的rpc请求最终将对backendsession的改变反馈到session.
///
/// 本测试中创建cobackendsession,并创建一个与session对应的backendsession，然后分别测试backendsessionservice的
/// 的方法.
package main

import (
	"component/cobackendsession"
	"component/corpcclient"
	"context"
	"fmt"
	"github.com/cihub/seelog"
)

func main() {
	ctx := context.GetContext()

	defer seelog.Flush()

	currentServer := make(map[string]interface{})
	currentServer["id"] = "chat-1"
	currentServer["serverType"] = "chat"
	currentServer["host"] = "127.0.0.1"
	currentServer["port"] = 8888

	ctx.CurrentServer = currentServer

	cobackendS := cobackendsession.NewCoBackendSession()
	corpcC := corpcclient.NewCoRpcClient()

	cobackendS.Start()
	corpcC.Start()

	backendS := cobackendS.GetBackendSessionBySID("connector-1", 1)

	if backendS != nil {
		fmt.Printf("--------------GetBackendSessionBySID(%v,%v)--------------\n", "connector-1", 1)
		fmt.Printf("frontendid:%v sid:%v  uid:%v\n", backendS.GetFrontendID(), backendS.GetID(), backendS.GetUID())
		fmt.Println("---------------------------------------------------------\n")

		backendS.UnbindUID(backendS.GetUID())
		fmt.Printf("--------------UnBindUID(%v)--------------\n", backendS.GetUID())
		backendSN := cobackendS.GetBackendSessionBySID("connector-1", 1)

		if backendSN != nil {
			fmt.Printf("--------------After UnBindUID(%v)--------------\n", backendS.GetUID())
			fmt.Printf("frontendid:%v sid:%v  uid:%v\n", backendSN.GetFrontendID(), backendSN.GetID(), backendSN.GetUID())
			fmt.Println("-----------------------------------------------------\n")
		}

		backendS.SetOpt("age", 12)
		fmt.Printf("--------------SetOpt(%v,%v)--------------\n", "age", 12)
		backendS.PushOpt("age")
		fmt.Printf("--------------PushOpt(%v)--------------\n", "age")
	}

	backendS = cobackendS.GetBackendSessionBySID("connector-1", 1)

	if backendS != nil {
		fmt.Printf("--------------GetBackendSessionBySID(%v,%v)--------------\n", "connector-1", backendS.GetID())
		fmt.Printf("frontendid:%v sid:%v  uid:%v,opts:%v\n", backendS.GetFrontendID(), backendS.GetID(), backendS.GetUID(), backendS.GetOpts())
		fmt.Println("--------------------------------------------------------------\n")
	}

	infoByUID := cobackendS.GetBackendSessionsByUID("connector-1", "Li Si")
	if infoByUID != nil {
		fmt.Printf("--------------GetBackendSessionSByUID(%v,%v)--------------\n", "connector-1", "Li Si")
		for _, elem := range infoByUID {
			fmt.Printf("frontend id:%v,id:%v,uid:%v \n", elem.GetFrontendID(), elem.GetID(), elem.GetUID())
		}
		fmt.Println("---------------------------------------------------\n")
	}

	cobackendS.KickBySID("connector-1", 1, "no reason")
	fmt.Printf("--------------KickBySID(%v,%v,%v)--------------\n", "connector-1", 1, "no reason")
	backendK := cobackendS.GetBackendSessionBySID("connector-1", 1)
	if backendK == nil {
		fmt.Printf("--------------After Kick sid<%v>,session not exists--------------\n\n", 1)
	}
}
