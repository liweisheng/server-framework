package main

import (
	"context"
	"fmt"

	"module"
	"os"
	"pomelo_admin"
	// "time"
)

func main() {
	Context := context.NewContext()
	Context.Env = "production"
	Context.MasterInfo["id"] = "master_server"
	Context.MasterInfo["host"] = "127.0.0.1"
	Context.MasterInfo["port"] = 3005
	// servers := make([]map[string]interface{}, 10)
	// singleServ := make(map[string]interface{})
	// singleServ["id"] = "connector"
	// singleServ["host"] = "127.0.0.1"
	// singleServ["port"] = 3150

	// append(servers)
	timer := &module.Timer{"timer", "request", 5}
	Context.RegisteModule(timer)
	fmt.Fprintf(os.Stdout, "Context.Modules: <%v>\n", Context.Modules)
	MC := pomelo_admin.NewMasterConsoleService(Context)
	MC.Start()
	// time.Sleep(5 * time.Second)
	select {
	case <-Context.Ch:
		fmt.Println("Master Server Exit... ")
		break
	}
}
