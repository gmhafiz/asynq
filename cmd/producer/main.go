package main

import (
	"tasks/internal/server"
)

const Version = "v0.1.0"

// Producer receives a new Unit of Work (UoW) from external services, and queues
// it to a Redis instance or Redis Cluster.
// You must send a unique key in the HTTP header X-Request-ID to ensure
// idempotent Unit of Work (UoW)
func main() {
	s := server.New(Version)
	s.Init()
	s.Run()
}
