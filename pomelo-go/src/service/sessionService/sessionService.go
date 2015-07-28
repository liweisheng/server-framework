/**
 *author:liweisheng date:2015/07/08
 */

/**
 *sessionService用于管理session.每一个session用来维护一条用户连接的所有信息.
 *
 *sessionService中定义两个struct：SessionService和Session.
 *其中Session维护一条用户连接，它维护的信息包括底层通信的socket，用户id等.
 *SessionService则负责管理所有的Session,同时也提供对Session方法的外层封装.
 */
package sessionService

import (
	"connector"
	"fmt"
	seelog "github.com/cihub/seelog"
	"os"
	// "context"
)

///session状态
const (
	SS_CLOSED = iota
	SS_INITED = iota
)

type SessionService struct {
	sessions  map[uint32]*Session    /// id->session
	uidMap    map[string][]*Session  /// uid->sessions
	multiBind bool                   /// true表示允许uid绑定多个session
	options   map[string]interface{} /// 选项
}

/// 创建新的SessionService.
///
/// @param opts 指定创建SessionService的选项，选项是name:value的形式
/// @return *SessionService
func NewSessionService(opts map[string]interface{}) *SessionService {
	if nil == opts {
		opts = make(map[string]interface{})
	}
	sessions := make(map[uint32]*Session)
	uidmap := make(map[string][]*Session)
	multibind := opts["multiBind"].(string) != ""
	return &SessionService{sessions, uidmap, multibind, opts}
}

/// 创建并返回新的session,同时SessionService内部维持此新的session
///
/// @param sid 每个session唯一的id，不应当为空,但是CreateSession内部不检查sid唯一性.
/// @param frontendID 创建此session前端服务器id.
/// @param socket 由session保持的底层连接,连接来自客户端
/// @return *Session
func (ss *SessionService) CreateSession(sid uint32, frontendID string, socket connector.Socket) *Session {
	session := newSession(sid, frontendID, socket)
	ss.sessions[sid] = session
	return session
}

/// 将用户id与session绑定.
///
/// @param uid 用户id,调用者确保不为空.
/// @param sid 待绑定的session的id.
/// @return 绑定过程出错则返回错误，绑定成功返回nil
func (ss *SessionService) BindUID(uid string, sid uint32) error {
	session := ss.sessions[sid]
	if nil == session {
		fmt.Fprintf(os.Stderr, "Error: Failed to find session with sid<%v>\n", sid)
		return fmt.Errorf("Error: Failed to find session with sid<%v>\n", sid)
	}

	if session.Uid != "" {
		if session.Uid == uid {
			//session已经绑定的相同的uid,返回
			return nil
		} else {
			fmt.Fprintf(os.Stderr, "Error: session with sid<%v> has already bound with uid<%v>\n", sid, uid)
			return fmt.Errorf("Error: session with sid<%v> has already bound with uid<%v>\n", sid, uid)
		}
	}

	sessions, ok := ss.uidMap[uid]
	if ss.multiBind == false && ok {
		//multiBind == false禁止同一个uid绑定到多个session
		fmt.Fprintf(os.Stderr, "Error: single uid can not be binded to multi sessions\n")
		return fmt.Errorf("Error: single uid can not be binded to multi sessions\n")
	}

	for _, elem := range sessions {
		if elem.Uid == uid {
			//已经有session绑定uid
			return nil
		}
	}
	session.bindUID(uid)
	ss.uidMap[uid] = append(ss.uidMap[uid], session)
	return nil

}

/// 解除用户id与session的绑定关系.
///
/// @param uid 用户id.
/// @param sid session id.
/// @return {error} 成功解除绑定放回nil，否者返回error.
func (ss *SessionService) UnbindUID(uid string, sid uint32) error {
	session, ok := ss.sessions[sid]
	if ok == false {
		//没有id为sid的session
		fmt.Fprintf(os.Stderr, "Error: Failed to find session with sid<%v>\n", sid)
		return fmt.Errorf("Error: Failed to find session with sid<%v>\n", sid)
	}

	if session.Uid != uid {
		fmt.Fprintf(os.Stderr, "Error: session with sid<%v> has not bind with uid<%v>\n", sid, uid)
		return fmt.Errorf("Error: session with sid<%v> has not bind with uid<%v>\n", sid, uid)
	}

	sessions, ok := ss.uidMap[uid]

	///将绑定uid的session从uid->session的map中移除
	if ok == true {
		var index int = -1
		for i, v := range sessions {
			if v.Id == sid {
				index = i
				break
			}
		}

		if index >= 0 {
			sessions = append(sessions[0:index], sessions[index+1:]...)
			if len(sessions) == 0 {
				delete(ss.uidMap, uid)
			} else {
				ss.uidMap[uid] = sessions
			}
		}
	}

	session.unbindUID()
	return nil

} //end UnBindUID

/// 通过session id返回相应的session.
///
/// @param sid session id
/// @param 有sid对应的session则返回否则返回nil.
func (ss *SessionService) GetSessionByID(sid uint32) *Session {
	return ss.sessions[sid]
}

/// 通过用户id返回用户绑定的所有session.
///
/// @param uid 用户id.
/// @return 返回nil或者用户绑定的session数组.
func (ss *SessionService) GetSessionsByUID(uid string) []*Session {
	return ss.uidMap[uid]
} //end GetSessionByUID

/// 通过session id来移除session，如果session有绑定uid，
/// 则同时从uid->sessions中移除session
///
/// @param sid 移除的session id
func (ss SessionService) RemoveSessionByID(sid uint32) {
	session, ok := ss.sessions[sid]

	if ok == true {
		uid := session.Uid
		delete(ss.sessions, sid)
		sessions := ss.uidMap[uid]

		if sessions != nil {
			var index int = -1

			for i, v := range sessions {
				if v.Id == sid {
					index = i
					break
				}
			}

			if index >= 0 {
				sessions = append(sessions[:index], sessions[index+1:]...)
			}

			if len(sessions) == 0 {
				delete(ss.uidMap, uid)
			} else {
				ss.uidMap[uid] = sessions
			}

		}
	}
} //end RemoveSessionByID

/// 断开所有uid绑定的连接.
///
/// @param uid 用户id.
/// TODO:目前只是强制断开连接，应该在断开连接之前发送发送提示信息.
/// XXX: 目前断开所有uid绑定的连接后并没有移除相应的session
func (ss *SessionService) KickByUID(uid string, reason string) {
	sessions := ss.GetSessionsByUID(uid)

	sess := make([]*Session, 0, 20)
	if sessions != nil {
		for _, elem := range sessions {
			sess = append(sess, elem)
		}

		for _, elem := range sess {
			ss.sessions[elem.Id].close(reason)
		}
	}
} //end KickByUID

/// 断开sid指定的session的连接.
///
/// @param sid session id.
/// TODO: 在断开连接之前应当发送一段由reason表示提示信息.
/// XXX: 断开连接后并未移除session.
func (ss *SessionService) KickBySessionID(sid uint32, reason string) {
	session, ok := ss.sessions[sid]

	if ok {
		session.close(reason)
	}
}

/// 通过session id获得客户端的地址信息.
///
/// @param sid session id
/// @return 返回地址信息或者nil.在成功返回的情况下，返回值的格式为name:value形式
///  如host:127.0.0.1 port:10000.
func (ss *SessionService) GetClientAddrBySID(sid uint32) map[string]string {
	session, ok := ss.sessions[sid]

	if ok == true {
		return session.Socket.RemoteAddress()
	}

	return nil
}

func (ss *SessionService) PushOpt(sid uint32, key string, value interface{}) {
	session, ok := ss.sessions[sid]

	if ok == true {

		seelog.Infof("Push opt with key<%v> value<%v> to session<%v>", key, value, sid)
		session.SetOpt(key, value)
	}
}

func (ss *SessionService) PushAllOpts(sid uint32, opts map[string]interface{}) {
	session, ok := ss.sessions[sid]

	if ok == true {
		seelog.Infof("Push all opts to session<%v>", sid)
		session.SetAllOpts(opts)
	}
}

/// 发送信息到session id绑定的连接上.
///
/// @param sid session id
/// @param msg 发送的信息,发送的信息应该有connector模块提供的encode函数编码过的信息
/// XXX: 应该使用日志而不是输出到终端.
func (ss *SessionService) SendMsgBySID(sid uint32, msg []byte) {
	session, ok := ss.sessions[sid]

	if ok == true {
		session.send(msg)
	} else {
		fmt.Fprintf(os.Stderr, "Try to send message to session with sid<%v>,which dose not exist\n", sid)
		return
	}
}

/// 发送信息给指定的用户id，如果一个用户id绑定多个session，则每个session都会收到信息.
///
/// @param uid  用户id.
/// @param msg 发送的消息,发送的信息应该有connector模块提供的encode函数编码过的信息
func (ss *SessionService) SendMsgByUID(uid string, msg []byte) {
	sessions, ok := ss.uidMap[uid]
	if ok == true {
		for _, elem := range sessions {
			elem.send(msg)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Try to send message to user with uid<%v>,which dose not exist\n", uid)
		return
	}

}

/// XXX: 应该使用日志记录发送失败.
func (ss *SessionService) SendBatchBySID(sid uint32, msgs ...[]byte) {
	session, ok := ss.sessions[sid]

	if ok == true {
		session.sendBatch(msgs...)
	} else {
		fmt.Fprintf(os.Stderr, "Try to send batch of messages to session with sid<%v>,which dose not exist\n", sid)
		return
	}
}

/// 批量发送信息给指定用户id.
///
/// @param uid 用户id
func (ss *SessionService) SendBatchByUID(uid string, msgs ...[]byte) {
	sessions, ok := ss.uidMap[uid]

	if ok == true {
		for _, elem := range sessions {
			elem.sendBatch(msgs...)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Try to send batch of messages to user with uid<%v>,which dose not exist\n", uid)
		return
	}
}

type Session struct {
	Status     int8
	Id         uint32
	Uid        string
	FrontendID string
	Socket     connector.Socket
	Opts       map[string]interface{}
}

/// 创建新的session,创建时的参数都不可以为空，由调用者保证.
func newSession(sid uint32, frontendId string, socket connector.Socket) *Session {
	opts := make(map[string]interface{})
	return &Session{SS_INITED, sid, "", frontendId, socket, opts}
}

/// 绑定用户ID.
///
/// @param uid {string} 用户id
func (s *Session) bindUID(uid string) {
	s.Uid = uid
}

///解除用户ID的和session的绑定.
func (s *Session) unbindUID() {
	s.Uid = ""
}

func (s *Session) SetOpt(key string, value interface{}) {
	s.Opts[key] = value
}

func (s *Session) GetOpt(key string) interface{} {
	return s.Opts[key]
}

func (s *Session) SetAllOpts(opts map[string]interface{}) {
	s.Opts = opts
}

func (s *Session) GetAllOpts() map[string]interface{} {
	return s.Opts
}

/// 关闭session，关闭session会关闭session保持的socket，
/// 同时将session信息从sessionService中移除.
///
/// TODO: 应该在断开连接之前给客户端发送一段简短提示信息.
func (s *Session) close(reason string) {
	if s.Status == SS_CLOSED {
		return
	}

	s.Status = SS_CLOSED
	// s.socket.RemoteAddress()
	s.Socket.Disconnect()
}

func (s *Session) send(msg []byte) {
	s.Socket.Send(msg)
}

func (s *Session) sendBatch(msgs ...[]byte) {
	s.Socket.SendBatch(msgs...)
}
