package cobackendsession

import (
	"context"
	"service/backendSessionService"
)

type CoBackendSession struct {
	*backendSessionService.BackendSessionService
}

func NewCoBackendSession() *CoBackendSession {
	ctx := context.GetContext()

	coBS, ok := ctx.GetComponent("cobackendsession").(*CoBackendSession)
	if ok == true {
		return coBS
	}
	bss := backendSessionService.NewBackendSessionService(ctx)
	coBS = &CoBackendSession{bss}

	ctx.RegisteComponent("cobackendsession", coBS)

	return coBS
}
