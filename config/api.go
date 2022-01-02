package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Api struct stores all API related configuration. The settings are parsed from
// environment variable. `default` struct tag ensures the field will be filled
// with a default value if the environment variable is empty.
type Api struct {
	Name       string `default:"queue"`
	Host       string `default:"localhost"`
	Port       int    `default:"4001"`
	RequestLog bool   `default:"false"`
}

func API() Api {
	var api Api
	envconfig.MustProcess("API", &api)

	return api
}
