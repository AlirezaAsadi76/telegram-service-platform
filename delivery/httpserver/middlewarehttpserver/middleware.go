package middlewarehttpserver

import "telegram-service-platform/service/authservice"

type Middleware struct {
	authSvc    *authservice.Service
	authConfig authservice.Config
}

func New(authSvc *authservice.Service, authConfig authservice.Config) Middleware {
	return Middleware{
		authSvc:    authSvc,
		authConfig: authConfig,
	}
}
