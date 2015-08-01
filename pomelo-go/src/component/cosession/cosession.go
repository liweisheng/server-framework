/**
 * author:liweisheng date:2015/07/15
 */

/**
 * cosesssion 是对sessionService的代理.
 */

package cosession

import (
	"context"
	"github.com/cihub/seelog"
	"service/sessionService"
)

type CoSession struct {
	*sessionService.SessionService
}

/// 创建CoSession.
func NewCoSession() *CoSession {
	ctx := context.GetContext()

	cosess, ok := ctx.GetComponent("cosession").(*CoSession)
	if ok == true {
		return cosess
	}
	ss := sessionService.NewSessionService(ctx.AllOpts["cosession"])

	cosess = &CoSession{ss}

	ctx.RegisteComponent("cosession", cosess)
	seelog.Info("CoSession create successufully")
	return cosess
}
