package main

import (
	"log"
	"tasks/task"

	"github.com/hibiken/asynq"

	"tasks/internal/server"
	"tasks/task/email"
)

const Version = "v0.1.0"

// Consumer dequeue tasks from Redis.
// You may have any number of consumer APIs (horizontal scaling) connecting to
// the same Redis cluster. Only one server will dequeue a task. Though you
// cannot choose which task goes to which server.
//
// One workaround is to customize the mux handler to only receive certain types
// of tasks that will only be handled in a particular machine. Those tasks may
// need a beefier RAM or CPU processing power. This means you will need to
// compile a different binary than the rest of your other consumer APIs and
// place them in that machine.
//
// These tasks are handled by registered Handlers in the route multiplexer (mux)
// below.
func main() {
	s := server.New(Version)
	s.Init()

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	// we inject database dependency into the handler
	mux.Handle(task.TypeEmailDelivery, email.NewEmailProcessor(s.DB()))
	// ...register other handlers...

	if err := s.AsynqServer().Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
