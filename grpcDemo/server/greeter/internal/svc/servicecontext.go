package svc

import "greeter/internal/config"

type ServiceContext struct {
	AtomicCfg *config.AtomicConfig
}

func NewServiceContext(ac *config.AtomicConfig) *ServiceContext {
	return &ServiceContext{
		AtomicCfg: ac,
	}
}

// Config returns the current active config (always fresh, safe for hot-reload).
func (sc *ServiceContext) Config() *config.Config {
	return sc.AtomicCfg.Load()
}
