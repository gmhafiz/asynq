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

// consumer dequeue tasks from Redis.
// You may have any number of consumer APIs (horizontal scaling) connecting to
// the same Redis cluster. Only one server will dequeue a task. Though you
// cannot choose which task goes to which server.
//
// These tasks are handled by registered Handlers in the route multiplexer (mux)
// below.
func main() {
	s := server.New(Version)
	s.Init()

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	// we inject database dependency into the handler
	mux.Handle(tasks.TypeEmailDelivery, email.NewEmailProcessor(s.DB()))
	// ...register other handlers...

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	go func() {
		if err := s.AsynqServer().Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	<-quit

	fmt.Println("...waiting for all tasks to complete")
}
