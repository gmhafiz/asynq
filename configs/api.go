package configs

import (
	"net"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Api struct {
	Name              string `default:"dribl_queue"`
	Host              net.IP
	Port              string
	ReadTimeout       time.Duration `default:"5s"`
	ReadHeaderTimeout time.Duration `default:"5s"`
	WriteTimeout      time.Duration `default:"10s"`
	IdleTimeout       time.Duration `default:"120s"`
	RequestLog        bool          `default:"false"`
}

func API() Api {
	var api Api
	envconfig.MustProcess("API", &api)

	return api
}
