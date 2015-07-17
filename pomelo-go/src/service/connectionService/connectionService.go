/**
 * author:liweisheng date:2015-07-13
 */

/**
 *connectionSevice维护连接的统计信息.
 */
package connectionService

type ConnectionService struct {
	serverID    string
	connCount   uint16
	loggedCount uint16
	loggedInfo  map[string]interface{}
}

func NewConnectionService(serverid string) *ConnectionService {
	loggedInfo := make(map[string]interface{})

	return &ConnectionService{serverid, 0, 0, loggedInfo}
}

/// 增加登录用户以及其信息.
///
/// @param uid 用户id
/// @param info 用户信息
/// @TODO 用户信息格式呈现没想好.
func (cs *ConnectionService) AddLoggedUser(uid string, info interface{}) {
	if _, ok := cs.loggedInfo[uid]; ok == false {
		cs.loggedCount += 1
	}

	cs.loggedInfo[uid] = info
}

/// 更新用户id对应的用户信息.
///
/// @param uid 需要更新的用户的用户id
/// @param info 用户的信息.
/// FIXIT: 直接将新的info赋值可能会导致旧的info中存在而新的info中不存在的信息丢失.
func (cs *ConnectionService) UpdateUserInfo(uid string, info interface{}) {
	if _, ok := cs.loggedInfo[uid]; ok == false {
		return
	}

	cs.loggedInfo[uid] = info
}

/// 将连接数加一.
func (cs *ConnectionService) IncreaseConnectionCount() {
	cs.connCount += 1
}

/// 根据用户id移除用户信息.
///
/// @param uid 用户id
func (cs *ConnectionService) RemoveLoggedUserByUID(uid string) {
	if _, ok := cs.loggedInfo[uid]; ok == true {
		cs.loggedCount -= 1
	}

	delete(cs.loggedInfo, uid)
}

/// 将当前连接数减1.
///
/// @param uid  如果uid不为nil,则同时移除uid指定的用户信息
func (cs *ConnectionService) DecreaseConnectionCount(uid string) {

	if cs.connCount > 0 {
		cs.connCount -= 1
	}

	if uid != "" {
		cs.RemoveLoggedUserByUID(uid)
	}
}

/// 获得统计信息.
///
/// 统计信息包括当前serverid，总连接数,总登录用户数,所有有登录用户信息,
/// 返回的信息的格式为name:value形式,如：
/// {"serverID":"connector-1","connectionCount":10,"loggedCount":6,"infos":[...]},
/// 其中[...]表示用户信息数组.
func (cs *ConnectionService) GetStatisticsInfo() map[string]interface{} {
	statistics := make(map[string]interface{})
	infos := make([]interface{}, 0)

	statistics["serverID"] = cs.serverID
	statistics["connectionCount"] = cs.connCount
	statistics["loggedCount"] = cs.loggedCount
	for _, elem := range cs.loggedInfo {
		infos = append(infos, elem)
	}

	statistics["infos"] = infos
	return statistics
}
