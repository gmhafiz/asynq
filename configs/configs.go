package configs

import (
	"fmt"

	"github.com/joho/godotenv"
)

// Configs is the main struct that holds all other configuration structs. It
// serves as the main entry to all other configuration types.
type Configs struct {
	Api      Api
	Database Database
	Redis    Redis
}

func New() *Configs {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	return &Configs{
		Api:      API(),
		Database: DataStore(),
		Redis:    NewRedis(),
	}
}
