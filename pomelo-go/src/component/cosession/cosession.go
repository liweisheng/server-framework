/**
 * author:liweisheng date:2015/07/15
 */

/**
 * cosesssion 是对sessionService的代理.
 */

package cosession

import (
	"context"
	"service/sessionService"
)

type CoSession struct {
	*sessionService.SessionService
	ctx *context.Context
}

/// 创建CoSession.
func NewCoSession(ctx *context.Context) *CoSession {
	ss := sessionService.NewSessionService(ctx.AllOpts["cosession"])

	cosess := &CoSession{ss, ctx}

	ctx.RegisteComponent("cosession", cosess)
	return cosess
}
