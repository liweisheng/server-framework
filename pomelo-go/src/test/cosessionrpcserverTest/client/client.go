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
		fmt.Printf("frontendif:%v sid:%v  uid:%v", backendS.GetFrontendID(), backendS.GetID(), backendS.GetUID())
	}

	ch := make(chan int)
	<-ch
}
