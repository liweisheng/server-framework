// master
package main

import (
	"context"
	"fmt"
	"pomelo_admin"
)

type Master struct {
	MasterInfo           map[string]interface{}
	Context              *context.Context
	MasterConsoleService *pomelo_admin.MasterConsoleService
}

///创建Master
func NewMaster(context *context.Context) Master {
	var masterInfo = make(map[string]interface{})
	var masterConsoleService = pomelo_admin.NewMasterConsoleService(context)
	var master = Master{masterInfo, context, &masterConsoleService}
	return master
}

func main() {
	fmt.Print("hello main")

}
