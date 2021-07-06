package configs

import (
	"net"

	"github.com/kelseyhightower/envconfig"
)

type Redis struct {
	Host      net.IP
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
