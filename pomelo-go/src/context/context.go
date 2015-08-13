// author:李为胜 2015-7-7

/*
package context用于记录一些全局的信息，包括：注册的modules，master服务器的配置信息，servers的配置信息
*/
package context

import (
	"fmt"
	seelog "github.com/cihub/seelog"
	"module"
	"os"
)

var globalContext *Context = nil

func init() {
	globalContext = newContext()
}

func GetContext() *Context {
	return globalContext
}

type Context struct {
	Ch                chan int8
	Env               string
	Server            interface{}
	Modules           []module.Module
	CurrentServer     map[string]interface{}
	MasterInfo        map[string]interface{}
	ServerInfo        map[string][]map[string]interface{}
	AllOpts           map[string]map[string]interface{} ///< 保存组件，模块的配置信息
	DefaultComponents map[string]interface{}
	Logger            seelog.LoggerInterface
}

//创建新的上下文
func newContext() *Context {
	ch := make(chan int8)
	mods := make([]module.Module, 0, 10)

	curS := make(map[string]interface{})
	masterInfo := make(map[string]interface{})
	serverInfo := make(map[string][]map[string]interface{})
	allOpts := make(map[string]map[string]interface{})
	defaultComponents := make(map[string]interface{})
	logger, err := seelog.LoggerFromConfigAsFile("./logConfig.xml")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Fail to create logger,error message:<%v>\n", err.Error())
		os.Exit(1)
	}
	if err := seelog.UseLogger(logger); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Fail to use logger,error message:<%v>\n", err.Error())
		os.Exit(1)
	}
	return &Context{ch, "", mods, nil, curS, masterInfo, serverInfo, allOpts, defaultComponents, logger}
}

/// 向上下文中注册一个module.
func (ctx *Context) RegisteModule(mod module.Module) {
	ctx.Modules = append(ctx.Modules, mod)
	fmt.Println(len(ctx.Modules))
}

func (ctx *Context) RegisteComponent(name string, comp interface{}) {
	seelog.Tracef("Registe component <%v>", name)
	ctx.DefaultComponents[name] = comp
}

/// 根据组件名称获得组件实例.
func (ctx *Context) GetComponent(name string) interface{} {
	return ctx.DefaultComponents[name]
}

// func (ctx *context) SetMasterInfo(mi map[string]interface{}) {
// 	ctx.MasterInfo = mi
// }

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Error Exit: %s\n", err.Error())
		os.Exit(1)
	}
}

func (ctx *Context) GetServerID() string {
	return ctx.CurrentServer["id"].(string)
}

func (ctx *Context) GetServerType() string {
	return ctx.CurrentServer["serverType"].(string)
}

func (ctx *Context) GetCurrentServerInfo() map[string]interface{} {
	return ctx.CurrentServer
}

/// 根据server id获得server的详细信息.
///
/// @param id server id
/// @return server的详细信息
/// XXX: 当前实现只是为了测试.
func (ctx *Context) GetServerInfoByID(id string) map[string]interface{} {
	return ctx.CurrentServer
}

/// 根据服务器类型获得server id,获得server id可以配置route规则，来从多个同类型服务器中
/// 选择要返回的server id.
/// @param stype 服务器类型.
/// @return server id
/// TODO: 未实现.
func (ctx *Context) GetServerIDByType(stype string) string {
	return ""
}
