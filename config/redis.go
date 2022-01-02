package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Redis struct {
	Host      string
	Port      string
	Addresses string
	Name      int
	User      string
	Pass      string
	CacheTime int
}

func NewRedis() Redis {
	var cache Redis
	envconfig.MustProcess("REDIS", &cache)

	return cache
}
