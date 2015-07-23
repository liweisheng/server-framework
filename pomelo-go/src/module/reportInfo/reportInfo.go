/**
 * author:liweisheng date:2015/07/14
 */

/**
 * 该模块实现的功能是：定期上报加载该模块的服务器的信息.
 * 加载模块的monitor服务定期的执行MonitorHandler上报自身服务的信息（serverType，serverID，ip，clientPort，port）
 * 该模块底层使用redis来上报自身的信息, 并设置信息过期时间为MonitorHandler执行周期的2倍.
 */
package reportInfo

import (
	"agent"
	"log"
	//"module"
	// "fmt"
	"net"
	"redis"
	"strconv"
	"sync"
)

var rwLock sync.RWMutex

type ReportInfo struct {
	moduleId   string
	moduleType string
	interval   int16
	client     redis.Client
}

func NewReportInfo(mid string, mt string, intv int16, reportAddr string, reportPort int) (*ReportInfo, error) {

	spec := redis.DefaultSpec().Host(reportAddr).Port(reportPort)
	cli, err := redis.NewSynchClientWithSpec(spec)

	if err != nil {
		return nil, err
	}

	return &ReportInfo{mid, mt, intv, cli}, nil
}

func (ri *ReportInfo) getServerInfoAsString(si map[string]interface{}) string {
	id := si["id"].(string)
	host := si["host"].(string)
	port := strconv.Itoa(si["port"].(int))
	clientPort := strconv.Itoa(si["clientPort"].(int))
	frontend := si["frontend"].(string)
	return id + ":" + host + ":" + port + ":" + clientPort + ":" + frontend
}

func (ri *ReportInfo) MonitorHandler(ag agent.Agent, conn net.Conn, msg map[string]interface{}) {
	log.Println("MonitorHandler is called")
	rwLock.Lock()
	defer rwLock.Unlock()
	// msg == nil表示是定期执行的任务.
	if msg == nil {
		if ag.IsMaster() {
			return
		}
		monitorAgent, ok := ag.(*agent.MonitorAgent)
		if ok == false {
			log.Fatal("MonitorHandler: Fail to convert to MonitorAgent")
		}
		serverInfo := monitorAgent.GetServerInfo()
		if serverInfo != nil {
			infoAsString := ri.getServerInfoAsString(serverInfo)
			ri.client.Sadd(serverInfo["serverType"].(string), []byte(serverInfo["id"].(string)))
			ri.client.Set(serverInfo["id"].(string), []byte(infoAsString))
			ri.client.Expire(serverInfo["id"].(string), int64(2*ri.interval))
		}

	}
}

func (ri *ReportInfo) MasterHandler(agent agent.Agent, conn net.Conn, msg map[string]interface{}) {}
func (ri *ReportInfo) ClientHandler(agent agent.Agent, conn net.Conn, msg map[string]interface{}) {}
func (ri *ReportInfo) Start() {
}
func (ri *ReportInfo) ModuleID() string {
	return ri.moduleId
}
func (ri *ReportInfo) GetType() string {
	return ri.moduleType
}
func (ri *ReportInfo) GetInterval() int16 {
	return ri.interval
}
