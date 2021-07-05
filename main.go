package main

import (
	"log"
	"tasks/tasks"
	"time"

	"github.com/hibiken/asynq"
)

const (
	redisAddr = "0.0.0.0:6379"
	TimeLayout = "2006-01-02 15:04"
)

func main() {
	// for single redis instance
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	// for connecting to the Redis cluster
	//client := asynq.NewClient(asynq.RedisClusterClientOpt{Addrs: []string{
	//	"0.0.0.0:7000", "0.0.0.0:7001", "0.0.0.0:7002",
	//}})
	defer client.Close()

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------

	task, err := tasks.NewEmailDeliveryTask(43, "some:template:id")

	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)


	// ------------------------------------------------------------
	// Example 2: Schedule task to be processed in the future.
	//            Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------

	//info, err = client.Enqueue(task, asynq.ProcessIn(24*time.Hour))
	myDate, err := time.Parse(TimeLayout, "2021-07-04 20:30")
	if err != nil {
		log.Fatalf("could not parse date: %v", err)
	}
	info, err = client.Enqueue(task, asynq.ProcessAt(myDate))
	if err != nil {
		log.Fatalf("could not schedule task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)


	// ----------------------------------------------------------------------------
	// Example 3: Set other options to tune task processing behavior.
	//            Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------

	client.SetDefaultOptions(tasks.TypeImageResize, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))

	task, err = tasks.NewImageResizeTask("https://example.com/myassets/image.jpg")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err = client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ---------------------------------------------------------------------------
	// Example 4: Pass options to tune task processing behavior at enqueue time.
	//            Options passed at enqueue time override default ones.
	// ---------------------------------------------------------------------------

	info, err = client.Enqueue(task, asynq.Queue("critical"), asynq.Timeout(30*time.Second))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}