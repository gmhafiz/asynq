// Package config stores all configurations loaded from .env file or
// environment variable into structs.
//
// Other files in this package defines its own configurations. This file glues
// all of them together with `Config` struct.
package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

// Config is the main struct that holds all other configuration structs. It
// serves as the main entry to all other configuration types.
type Config struct {
	Api      Api
	Database Database
	Redis    Redis
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	return &Config{
		Api:      API(),
		Database: DataStore(),
		Redis:    NewRedis(),
	}
}
