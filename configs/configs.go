package configs

import (
	"fmt"

	"github.com/joho/godotenv"
)

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
