package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

var reqID uint32 = 0

func send(route string, body []byte, conn net.Conn) {
	msg := make(map[string]interface{})
	msg["reqID"] = fmt.Sprintf("%v", reqID)
	msg["route"] = route
	msg["body"] = string(body)
	fmt.Println("in Send")
	sendBuf, err := json.Marshal(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Json Marshall msg<%v> error,error message<%v>\n", msg, err.Error())
		os.Exit(1)
	}
	rst := make(map[string]interface{})
	if err = json.Unmarshal(sendBuf, &rst); err == nil {
		fmt.Printf("Json UnMarshall<%v> -> <%v>\n", string(sendBuf), rst)
	} else {
		fmt.Printf("Json UnMarshall<%v> error,error message<%v>\n", string(sendBuf), err.Error())
		os.Exit(1)
	}
	bufLen := uint16(len(sendBuf))
	lenBuf := make([]byte, 2)

	lenBuf[0] = byte((bufLen >> 8) & 0xFF)
	lenBuf[1] = byte(bufLen & 0xFF)
	_, err1 := conn.Write(lenBuf)

	if err1 == nil {
		fmt.Printf("msg len:high<%v> low<%v>\n", lenBuf[0], lenBuf[1])
	} else {
		fmt.Printf("send lenBuf error,error message<%v>\n", err1)
		os.Exit(1)
	}

	//	_, err2 := conn.Write(sendBuf)
	_, err2 := conn.Write([]byte("fuck you"))
	if err2 == nil {
		fmt.Printf("send msg<%v>,send len<%v>\n", string(sendBuf), len(sendBuf))
	} else {
		fmt.Printf("send sendBuf error,error message<%v>\n", err2)
		os.Exit(1)
	}
	reqID++
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "too few args,args form: <host port>\n")
		os.Exit(1)
	}
	host := os.Args[1]
	port, err1 := strconv.Atoi(os.Args[2])
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "invalid port,need integer type,your input port: <port>\n", os.Args[2])
		os.Exit(1)
	}

	hostAndPort := fmt.Sprintf("%v:%v", host, port)

	conn, err2 := net.Dial("tcp", hostAndPort)
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Dial <%v> error,error message<%v>\n", hostAndPort, err2.Error())
		os.Exit(1)
	}
	//	defer conn.Close()

	//	route := "chat.Chatroom.Addin"
	msgBody := make(map[string]interface{})
	msgBody["type"] = "handshake"
	msgBody["body"] = "Hello,I am your father"
	body, err3 := json.Marshal(msgBody)
	if err3 != nil {
		fmt.Fprintf(os.Stderr, "Json Marshall <%v> error,error message<%v>\n", msgBody, err3.Error())
		os.Exit(1)
	}

	fmt.Printf("len: %v  text: %v\n", len(body), string(body))
	conn.Write([]byte("Fuck You!!!"))
	//	send(route, body, conn)
	ch := make(chan int)
	<-ch
}
