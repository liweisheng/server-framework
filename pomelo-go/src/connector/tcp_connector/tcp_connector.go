/*
author:liweisheng date:2015/07/08
*/

/*
实现tcp connector
*/
package tcp_connector

import (
	// "context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

var curID int32 = 0

type TcpConnector struct {
	host           string
	port           string
	opts           map[string]string
	registedEvents map[string]func(args ...interface{})
}

/// 创建新的TcpConnector
func NewTcpConnector(host string, port string, opts map[string]string) *TcpConnector {
	regE := make(map[string]func(args ...interface{}))
	return &TcpConnector{host, port, opts, regE}
}

/// 处理新接收到的连接.
///
/// 接收tcpSkt上的数据，并解析数据包，调用注册的message事件（函数回调）处理收到的数据.
///
/// @param tcpSkt 与客户端连接的socket.
func (tc *TcpConnector) HandleNewConnection(tcpSkt *TcpSocket) {
	const BUFSIZE uint16 = 1024 * 8
	var buff, recvBuff []byte
	var begin, end, packSize, unProcess uint16

	recvBuff = make([]byte, BUFSIZE)
	buff = make([]byte, BUFSIZE)
	for {

		n, err := tcpSkt.Receive(recvBuff)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Read from conn , err message :%v\n", err.Error())
			errEv, ok := tc.registedEvents["error"]
			if ok == false {
				fmt.Fprintf(os.Stderr, "Error: Can not find registed event handler<'error'>")
				os.Exit(1)
			}

			errEv(tcpSkt)
			break
		}
		if begin >= end {
			begin = 0
			end = 0
		}

		buff = append(buff[begin:end], recvBuff[0:n]...)
		// fmt.Fprintf(os.Stdout, "current buff is: %v\n", buff[:])
		unProcess = uint16(len(buff))
		begin = 0

		for unProcess >= 1 {
			packSize = uint16(0x00FF&buff[begin])<<8 + uint16(0x00FF&buff[begin+1])
			fmt.Fprintf(os.Stdout, "packsize is %v\n", packSize)
			if unProcess >= packSize {
				msg, err := tc.Decode(buff[begin+2 : begin+packSize])
				if err == nil {
					// goto DecodeErr
					msgEv, ok := tc.registedEvents["message"]
					if ok == false {
						fmt.Fprintf(os.Stderr, "Error: Can not find registed event handler<'message'>")
						os.Exit(1)
					}
					//处理接收到的消息
					go msgEv(msg)
				}
				unProcess -= packSize
				begin += packSize
			} else {
				break
			}
		} //end inner for
	} //end outter for
}

/// 监听服务器端口，接收新的连接.对于新来的连接首先调用为其注册的connection事件(函数回调)
/// 之后开始监听新的连接.
func (tc *TcpConnector) Start() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", tc.host+":"+tc.port)

	if err != nil {
		log.Fatal(err.Error())
	}
	// context.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	go func(ln *net.TCPListener) {
		for {
			conn, err := ln.AcceptTCP()
			// context.CheckError(err)
			if err != nil {
				log.Printf("AcceptTCP Error: error message<%v>\n", err.Error())
				// fmt.Fprintf(os.Stderr)
				os.Exit(0)
				// log.Fatalln(err.Error())
			}
			cb, ok := tc.registedEvents["connection"]
			if ok == false {
				log.Printf("Error: Fail to load <Events:'connection'>\n")
				// fmt.Fprintf(os.Stdout,)
				os.Exit(0)
			}
			tcpSocket := NewTcpSocket(curID, conn)
			cb(tcpSocket)
			go tc.HandleNewConnection(tcpSocket)
		} //end for
	}(listener)
} //end Start()

/// 为TcpConnector注册事件相应回调.
///
/// 至少要注册两个回调，一个是connection回调，用于新的连接到来时回调,
/// 一个是message回调，当有新的message到来时回掉.
///
/// @param evName 事件名称，connection事件名称"connection",message事件名称为"message".
/// @param callback 对应事件的回调函数.
func (tc *TcpConnector) RegistEvents(evName string, callback func(...interface{})) {
	tc.registedEvents[evName] = callback
}

/// Start之前调用，可以做一些必要的检查.
func (tc *TcpConnector) beforeStart() {
	if _, ok := tc.registedEvents["connection"]; ok == false {
		log.Fatalln("Did not find call-back function for 'connection' events")
	}

	if _, ok := tc.registedEvents["message"]; ok == false {
		log.Fatalln("Did not find call-back function for 'message' events")
	}
}

/// 解码收到的客户端信息.
///
/// 客户端的信息应该符合约定好的消息格式，否则解码失败.
///
/// @param buff 收到的信息.
/// @return rst {map[string]string} error 为nil表示解码成功,rst中以name:value形式.
func (tc *TcpConnector) Decode(buff []byte) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal(buff, &result)
	return result, err
}

/// 编码消息,编码成json格式.
///
/// @reqID 请求ID
/// @route 请求路由
/// @body 消息内容
/// @return result 第二返回值error为nil的情况下result中存放编码成json后的内容
func (tc *TcpConnector) Encode(reqID string, route string, body string) ([]byte, error) {
	msg := make(map[string]interface{})
	msg["id"] = reqID
	msg["route"] = route
	msg["body"] = body
	msgJson, err := json.Marshal(msg)

	if err != nil {
		return nil, err
	}
	return msgJson, nil
}
