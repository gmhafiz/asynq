package configs

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Api struct stores all API related configuration. The settings are parsed from
// environment variable. `default` struct tag ensures the field will be filled
// with a default value if the environment variable is empty.
type Api struct {
	Name              string        `default:"queue"`
	Host              string        `default:"localhost"`
	Port              int           `default:"4001"`
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
