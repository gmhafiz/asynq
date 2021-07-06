package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"

	"tasks/internal/server"
	"tasks/tasks"
	"tasks/tasks/email"
)

const Version = "v0.1.0"

func main() {
	s := server.New(Version)
	s.Init()

	srv := asynq.NewServer(
		// for single redis instance
		asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", s.Config().Redis.Addresses, s.Config().Redis.Port)},
		// for connecting to the Redis cluster
		//asynq.RedisClusterClientOpt{
		//	Addrs: []string{"0.0.0.0:7000", "0.0.0.0:7001", "0.0.0.0:7002"},
		//},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmailDelivery, email.HandleEmailDeliveryTask)
	mux.Handle(tasks.TypeImageResize, email.NewImageProcessor())
	// ...register other handlers...

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	<-quit

	fmt.Println("...waiting for all tasks to complete")
}
