package main

import (
	// "fmt"
	"context"
	"log"
	"os"
	"rpcserver"
	"strconv"
)

type Echo int

func (e *Echo) Hi(arg string, reply *string) error {
	*reply = "Echo " + arg
	return nil
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("too few args to server,commandline form <command host port>")
	}

	context.NewContext()
	host := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Fail to convert os.args[2] to port(type int),error message:", err.Error())
	}

	server := rpcserver.NewRpcServer(host, port)
	server.RegisteService(new(Echo))

	server.Start()
}
