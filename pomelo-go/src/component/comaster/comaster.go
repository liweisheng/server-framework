package comaster

import (
	"context"
	"github.com/cihub/seelog"
	"pomelo_admin"
)

type CoMaster struct {
	*pomelo_admin.MasterConsoleService
}

func NewCoMaster() *CoMaster {
	ctx := context.GetContext()
	coMaster, ok := ctx.GetComponent("comaster").(*CoMaster)
	if ok == true {
		return coMaster
	}

	mcs := pomelo_admin.NewMasterConsoleService(ctx)
	coMaster = &CoMaster{mcs}
	seelog.Infof("<%v> component CoMaster created...", ctx.GetServerID())
	ctx.RegisteComponent("comaster", coMaster)
	return coMaster
}
