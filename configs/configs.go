package configs

import (
	"log"

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
		log.Fatal(err)
	}

	return &Configs{
		Api:      API(),
		Database: DataStore(),
		Redis:    NewRedis(),
	}
}
