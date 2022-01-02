package queue

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
)

// Enqueue is a convenience function that defaults to using a supplied UUID
// for uniqueness and Retention() that keeps the completed task in the queue
// for a certain period. Keeping it in th queue prevents another task with the
// same task ID (UUID) from starting - this is useful for deduplication
// in the event of retries of same tasks.
func Enqueue(ctx context.Context, client *asynq.Client, task *asynq.Task, uuid string, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	var variadic []asynq.Option
	variadic = append(variadic, asynq.TaskID(uuid))
	variadic = append(variadic, asynq.Retention(5*time.Minute))
	variadic = append(variadic, opts...)

	return client.EnqueueContext(ctx, task, variadic...)
}
