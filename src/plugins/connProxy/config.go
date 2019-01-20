package connProxy

import "github.com/agui2200/roadrunner/service"

type ConnConfig struct {
	Conn map[string]struct {
		Host       string
		Port       uint
		Protocol   string
		PoolNum    uint
		ListenAddr string
	}
	Enable bool
}

func (c *ConnConfig) Hydrate(cfg service.Config) error {
	return cfg.Unmarshal(&c)
}
