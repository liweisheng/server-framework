package pomelo_admin

import (
	"agent"
	"context"
	"encoding/json"
	"fmt"
	"module"
	"net"
	"os"
)

const (
	SV_INIT  = iota
	SV_START = iota
	SV_CLOSE = iota
)

type packageInfo struct {
	ModuleID string
	Type     string
	Source   string
	Msg      map[string]interface{}
}

type ConsoleService interface {
	GetModuleByID(string) module.Module
	// HandleNewAcception(conn net.Conn)
	// Listen(hostAndPort string)
	Start()
	Stop()
	GetAgent() agent.Agent
}

///处理接受到的包，解析包中包含的信息，来决定调用的module
func processPackage(service ConsoleService, pack []byte, conn net.Conn) {
	fmt.Fprintf(os.Stdout, "Info: In processPackage...")
	var packInfo packageInfo
	var mod module.Module
	err := json.Unmarshal(pack, &packInfo)
	context.CheckError(err)
	mod = service.GetModuleByID(packInfo.ModuleID)
	if nil == mod {
		fmt.Fprintf(os.Stderr, "Error: Fail to get ModuleID: <%s>\n", packInfo.ModuleID)
		os.Exit(0)
	}

	fmt.Println("*************packInfo************")
	fmt.Println("ModuleID:", packInfo.ModuleID)
	fmt.Println("Type:", packInfo.Type)
	fmt.Println("Source:", packInfo.Source)
	fmt.Println("Msg:", packInfo.Msg)
	fmt.Println("*************packInfo************")
	switch packInfo.Source {
	case "master":
		mod.MasterHandler(service.GetAgent(), conn, packInfo.Msg)
	case "monitor":
		mod.MonitorHandler(service.GetAgent(), conn, packInfo.Msg)
	case "client":
		mod.ClientHandler(service.GetAgent(), conn, packInfo.Msg)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown Source:<%s>", packInfo.Source)
		return
	}
}

///处理连接上接收的数据
///
///BUG: 传递给processPackage参数有问题
func handlerConnectionRecv(service ConsoleService, conn net.Conn) {
	fmt.Fprintf(os.Stdout, "In handlerConnectionRecv, conn<%v>\n", conn)
	const BUFSIZE int16 = 512

	var recvBuff []byte
	var buff []byte
	var begin, end, packSize, unProcess int16

	recvBuff = make([]byte, BUFSIZE)
	buff = make([]byte, BUFSIZE)

	defer conn.Close()
	for {

		n, err := conn.Read(recvBuff)

		if err != nil {
			masterServ, ok := service.(*MasterConsoleService)
			if ok {
				masterServ.GetAgent().(*agent.MasterAgent).RemoveConnection(conn)
			}

			fmt.Fprintf(os.Stderr, "Error: Read from conn , err message :%v\n", err.Error())
			break
		}
		// fmt.Fprintf(os.Stdout, "what I receive is :%v", recvBuff[0:n])
		context.CheckError(err)
		if begin >= end {
			begin = 0
			end = 0
		}

		buff = append(buff[begin:end], recvBuff[0:n]...)
		// fmt.Fprintf(os.Stdout, "current buff is: %v\n", buff[:])
		unProcess = int16(len(buff))
		begin = 0

		// fmt.Fprintf(os.Stdout, "after read begin: %v,end: %v\n", begin, end)
		//		var reciveSize int16 = (end - begin + SIZE) % SIZE
		//		fmt.Fprintf(os.Stdout, "reciveSize %v\n", reciveSize)
		for unProcess >= 1 {
			packSize = int16(buff[begin])
			fmt.Fprintf(os.Stdout, "packsize is %v\n", packSize)
			if unProcess >= packSize {
				// for i := 1; int16(i) < packsize; i++ {
				// 	fmt.Fprintf(os.Stdout, "pos is %v\n", begin+int16(i))
				// 	fmt.Fprintf(os.Stdout, "recive %dth is %d\n", int16(i), buff[begin+int16(i)])
				// }
				fmt.Fprintf(os.Stdout, "Info: Packinfo :%v\n", buff[begin+1:begin+packSize])
				go processPackage(service, buff[begin+1:begin+packSize], conn)
				unProcess -= packSize
				begin += packSize
			} else {
				break
			}
			// fmt.Fprintf(os.Stdout, "after pocess begin: %v,end: %v\n", begin, end)
		} //end inner for
		//		fmt.Fprintf(os.Stdout, "after pocess begin: %v,end: %v\n", begin, end)
	} //end outter for

} //end func handleConnectionRecv
