/**
 * author:liweisheng date:2015/07/28
 */

/**
 * 包channelService包含Channel和ChannelServevice,Channel可以理解为一个广播组,
 * Channel包含多个session,每个session可以当作一个广播组的一个成员,channelService
 * 管理所有Channel.
 */
package channelService

import (
	"component/cochannelrpcserver"
	"component/corpcclient"
	"context"
	seelog "github.com/cihub/seelog"
)

type ChannelService struct {
	ctx         *context.Context
	coRpcClient *corpcclient.CoRpcClient
	channels    map[string]*Channel
}

/// 创建新的ChannelService,ChannelService管理用户创建的所有Channel,并代理用户对Channel的操作,
/// 用户发送给一个Channel里所有session的消息会通过rpc调用发送给对应的前端服务器，并通过前端服务
/// 器推送给客户端.
func NewChannelService() *ChannelService {
	ctx := context.GetContext()

	coRpcClient, ok := ctx.GetComponent("corpcclient").(*corpcclient.CoRpcClient)

	if ok == false {
		coRpcClient = corpcclient.NewCoRpcClient()
	}

	channels := make(map[string]*Channel)

	return &ChannelService{ctx, coRpcClient, channels}
}

/// 以name为名称创建新的Channel,如果指定名称的channel存在就返回存在的channel,否则创建新的.
///
/// @param name channel名称
/// @return {*Channel} 或者nil
func (cs *ChannelService) NewChannel(name string) *Channel {
	if name == "" {
		seelog.Error("NewChannel need a non-empty name")
		return nil
	}

	if channel, ok := cs.channels[name]; ok == true {
		seelog.Infof("channel with name<%v> already exists", name)
		return channel
	} else {
		seelog.Infof("channel with name<%v> doesn't exists,create new", name)
		channel := newChannel()
		cs.channels[name] = channel
		return channel
	}
}

/// 返回名称为name的channel，如果参数create为true，则指定名称的channel不存在时以name创建新的channel，
/// 否则返回nil.
///
/// @param name channel名称.
/// @param create 为true时名称为name的channel不存在是创建新的channel
/// @return *Channel 或者nil
func (cs *ChannelService) GetChannel(name string, create bool) *Channel {
	if channel, ok := cs.channels[name]; ok == true {
		return channel
	} else {
		if create == true {

			seelog.Infof("Channel with name<%v> doesn't exists ,create new", name)
			return cs.NewChannel(name)
		}

		return nil
	}

}

/// 销毁名称为name的channel
///
/// @param name 待销毁channel的名称
func (cs *ChannelService) DestroyChannel(name string) {
	seelog.Infof("Destroy channel<%v>", name)
	delete(cs.channels, name)
}

func (cs *ChannelService) PushMsgByUIDs(route string, msg map[string]interface{}, uid2servid map[string][]string) {

	// uidsX := make([]string, 0)

	// for uid, serverid := range uids {
	// 	if serverid == "" {
	// 		uidsX := append(uidsX, uid)
	// 	}
	// }

	// for uid := range uidsX {
	// 	seelog.Warnf("uid<%v> with empty serverid,just ignore it", uid)
	// 	delete(uids, uid)
	// }
	//serverid ->uids
	sid2uid := make(map[string][]string)

	for key, value := range uid2servid {
		for _, servid := range value {
			if servid == "" {
				seelog.Warnf("uid<%v> with empty serverid, just ignore it", key)
				continue
			} else {
				if uids, ok := sid2uid[servid]; ok == false {
					uids := make([]string, 5)
					uids = append(uids, key)
					sid2uid[servid] = uids
				} else {
					uids = append(uids, key)
					sid2uid[servid] = uids
				}
			}
		}
	}
	cs.sendMsgByGroup(route, msg, sid2uid)
}

/// 发送消息到groups中的所有用户.
///
/// @route 路由
/// @param msg 发送的消息
/// @param groups frontendID-> uids的映射
func (cs *ChannelService) sendMsgByGroup(route string, msg map[string]interface{}, groups map[string][]string) {
	seelog.Debugf("<%v> channelService sendMsgByGroup with route<%v> msg<%v> groups<%v>", cs.ctx.GetServerID(), route, msg, groups)

	args := make(map[string]interface{})

	args["route"] = route
	args["msg"] = msg

	coChanRpcS, ok := context.GetContext().GetComponent("cochannelrpcserver").(*cochannelrpcserver.CoChannelRpcServer)
	if ok == false {
		seelog.Error("Fail to get CoChannelRpcServer")
		return
	}

	sendMsg := func(serverid string) {
		if cs.ctx.GetServerID() == serverid {
			//TODO:直接调用channelRpcServer相关方法,无需rpc
			if err := coChanRpcS.PushMessage(args, nil); err != nil {
				seelog.Errorf("Fail to invoke PushMessage to send msg,error<%v>", err.Error())
				return
			}

		} else {
			//TODO:发起rpc调用.
			if err := cs.coRpcClient.RpcCall(serverid, "ChannelRpcServer.PushMessage", args, nil); err != nil {
				seelog.Errorf("<%v> fail to invoke rpc PushMessage error<%v>", context.GetContext().GetServerID(), err.Error())
				return
			}
		}
	}

	for frontendID, uids := range groups {
		//TODO:逐个发送msg.
		seelog.Infof("<%v> PushMessage to uids<%v> connected to frontend<%v>", context.GetContext().GetServerID(), uids, frontendID)
		args["uid"] = uids
		sendMsg(frontendID)
	}
}

type Channel struct{}

func newChannel() *Channel {
	return nil
}
