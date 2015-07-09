package pomelo_go

import (
	"context"
	//	"bufio"
	"encoding/json"
	"fmt"
	// "io"
	"os"
)

type pomelo_go struct {
	Context *context.Context
}

func NewPomeloGo() *pomelo_go {
	var ctx *context.Context = context.NewContext()

	return &pomelo_go{ctx}
}

///从./config/master.json中加载master信息
func (p *pomelo_go) LoadMasterInfo() {
	inputFile, err := os.Open("./config/master.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to open %v\n", "./config/master.json")
		os.Exit(1)
	}

	decoder := json.NewDecoder(inputFile)

	var masterInfo map[string]interface{}

	for {
		if err := decoder.Decode(&masterInfo); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to decode %v,Error Message:<%v>\n", "./config/master.json", err.Error())
			os.Exit(1)
		}
	}

	p.Context.MasterInfo = masterInfo[p.Context.Env]
}

///从./config/servers.json中加载server信息
func (p *pomelo_go) LoadServerInfo() {
	inputFile, err := os.Open("./config/servers.json")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to open %v\n", "./config/servers.json")
		os.Exit(1)
	}

	decoder := json.NewDecoder(inputFile)

	var serversInfo map[string]interface{}

	for {
		if err := decoder.Decode(&serversInfo); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to decode %v,Error Message:<%v>\n", "./config/servers.json", err.Error())
			os.Exit(1)
		}
	}

	p.Context.ServerInfo = serversInfo[p.Context.Env]
}

func (p *pomelo_go) SetEnv(env string) {
	p.Context.Env = env
}


///分析命令行参数
func (p *pomelo_go) ParseArgs{

}
