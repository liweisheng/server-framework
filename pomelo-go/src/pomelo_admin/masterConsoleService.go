package pomelo_admin

import (
	"agent"
	"context"
	// "encoding/json"
	"fmt"
	"module"
	"net"
	"os"
	"strconv"
)

/// MasterConsoleService实现consoleService接口
type MasterConsoleService struct {
	Context   *context.Context
	ModuleMap map[string]module.Module
	MAgent    *agent.MasterAgent
	Listner   *net.TCPListener
	status    int8
}

/// 创建新的MasterConsoleService.
///
/// 创建MasterConsoleService时首先会拿到所有注册在Context中的module，建立moduleID到
/// module的映射
func NewMasterConsoleService(ctx *context.Context) *MasterConsoleService {
	var modules = ctx.Modules
	var moduleMap = make(map[string]module.Module, 10)
	var mAgent = agent.NewMasterAgent(ctx.Ch, ctx.MasterInfo)
	for _, v := range modules {
		if v != nil {
			fmt.Fprintf(os.Stdout, "Registe module<%v> to master \n", v.ModuleID())
			moduleMap[v.ModuleID()] = v
		}

	}
	return &MasterConsoleService{ctx, moduleMap, mAgent, nil, SV_INIT}
}

///通过moduleID获取module
func (m *MasterConsoleService) GetModuleByID(moduleID string) module.Module {
	return m.ModuleMap[moduleID]
}

///处理接受到的新连接
func (m *MasterConsoleService) HandleNewAcception(conn net.Conn) {
	m.MAgent.AddConnection(conn)
	fmt.Fprintf(os.Stdout, "In HandlerNewAcception, conn<%v>\n", conn)
	defer conn.Close()
	handlerConnectionRecv(m, conn)
}

/// 开启监听端口，接收monitor的连接.
///
///hostAndPort   host:port
func (m *MasterConsoleService) Listen(hostAndPort string) {
	// fmt.Println("In  Listen")
	// fmt.Fprintf(os.Stdout, "tcp4 addr: %v\n", hostAndPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", hostAndPort)
	// fmt.Println("After Resolve")
	context.CheckError(err)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	// fmt.Fprintf(os.Stdout, "tcp4 addr: %v\n", tcpAddr)
	// fmt.Println("After Listen")
	context.CheckError(err)
	m.Listner = listener
	defer listener.Close()
	fmt.Fprintf(os.Stdout, "Info: Listening <%v>\n", hostAndPort)
	for {
		conn, err := listener.Accept()
		fmt.Fprintf(os.Stdout, "Info: New acception<%v>\n", conn.RemoteAddr())
		if nil != err {
			break

		} else {
			go m.HandleNewAcception(conn)
		}
	}
}

/// 启动MasterConsoleService.
///
/// 启动时，首先开启监听，遍历所有挂载的模块，如果模块配置类型为pull或者push,这开启定时
/// 调度，默认调度时间为5秒. 然后启动所有module.
func (m *MasterConsoleService) Start() {
	if m.status == SV_START {
		fmt.Println("Info: Master Server is started already")
		return
	}
	// for _, v := range m.Context.Modules {
	// 	if
	// 	m.ModuleMap[v.ModuleID()] = v
	// }

	host, ok := m.Context.MasterInfo["host"]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: Master host not set\n")
		os.Exit(1)
	}

	port, ok := m.Context.MasterInfo["port"]

	if !ok {
		fmt.Fprintf(os.Stderr, "Error: Master port not set\n")
		os.Exit(1)
	}

	host_s := host.(string)
	fmt.Println(port)
	port_i := port.(int)

	hostAndPort := host_s + ":" + strconv.Itoa(port_i)
	go m.Listen(hostAndPort)

	for _, v := range m.Context.Modules {
		if v != nil {
			if v.GetType() == "pull" {
				interval := v.GetInterval()
				if interval == 0 {
					interval = 5
				}
				go module.PeriodicScheduler(v.MasterHandler, m.MAgent, interval)
			}
		}
	}

	for _, v := range m.Context.Modules {
		if v != nil {
			v.Start()
		}
	}
}

func (m *MasterConsoleService) Stop() {
	if m.status != SV_START {
		return
	}

	m.Listner.Close()
	m.status = SV_CLOSE
}

func (m *MasterConsoleService) GetAgent() agent.Agent {
	return m.MAgent
}
