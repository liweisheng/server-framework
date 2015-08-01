package cochannel

import (
	"context"
	"github.com/cihub/seelog"
	"service/channelService"
)

type CoChannel struct {
	*channelService.ChannelService
}

func NewCoChannel() *CoChannel {
	ctx := context.GetContext()

	coChan, ok := ctx.GetComponent("cochannel").(*CoChannel)
	if ok == true {
		return coChan
	}
	chanS := channelService.NewChannelService()
	coChan = &CoChannel{chanS}
	ctx.RegisteComponent("cochannel", coChan)
	seelog.Info("CoChannel create successfully")
	return coChan

}
