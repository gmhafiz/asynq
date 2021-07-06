package main

import (
	"log"

	"tasks/internal/server"
)

const Version = "v0.1.0"

func main() {
	s := server.New(Version)
	s.Init()

	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
}

