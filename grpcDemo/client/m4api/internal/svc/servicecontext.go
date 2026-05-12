package svc

import (
	"client/greeterclient"
	"client/m4api/internal/config"
)

type ServiceContext struct {
	Config     config.Config
	GreeterRpc greeterclient.Greeter
}

func NewServiceContext(c config.Config, greeterRpc greeterclient.Greeter) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		GreeterRpc: greeterRpc,
	}
}
