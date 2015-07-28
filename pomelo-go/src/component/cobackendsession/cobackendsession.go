package cobackendsession

import (
	"context"
	"service/backendSessionService"
)

type CoBackendSession struct {
	*backendSessionService
}

func NewCoBackendSession() *CoBackendSession {
	ctx := context.GetContext()

	bss := backendSessionService.NewBackendSessionService(ctx)
	cobs := &CoBackendSession{bss}

	ctx.RegisteComponent("cobackendsession", cobs)
	return cobs
}
