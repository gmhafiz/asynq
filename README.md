An example of [asynq](https://github.com/hibiken/asynq) producer that creates tasks to be sent to redis

# Requirement

Go version >= 1.13

# Run

    go run main.go

# Send Signal

Find the process ID of this API by port

    lsof -t -i:4001

Send SIGTERM signal

    kill -SIGTERM <pid number>

Or one-liner

    kill -SIGTERM $(lsof -t -i:4001)