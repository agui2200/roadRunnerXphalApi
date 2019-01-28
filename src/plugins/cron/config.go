package cron

import "github.com/agui2200/roadrunner/service"

type Config struct {
	WorkDir string
}

// Hydrate must populate Config values using given Config source. Must return error if Config is not valid.
func (c *Config) Hydrate(cfg service.Config) error {
	return cfg.Unmarshal(&c)
}
