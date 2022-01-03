An example of [asynq](https://github.com/hibiken/asynq) producer and consumer 
that creates tasks to be sent to redis and process tasks from Redis queue. There
are two distinct microservices, consumer and producer, which allows you to scale
them independently.

# Table of Contents

- [Requirements](#requirements)
   * [Redis](#redis)
   * [(Optional) Protobuf](#-optional--protobuf)
   * [(Optional) Linting](#-optional--linting)
- [Run](#run)
   * [Database](#database)
   * [Producer](#producer)
   * [Consumer](#consumer)
- [Gracefully Shutdown](#gracefully-shutdown)
- [Usage](#usage)
   * [Producer](#producer-1)
   * [Consumer](#consumer-1)
- [Monitor](#monitor)
   * [asynqmon](#asynqmon)
- [Create a New Task](#create-a-new-task)
- [Appendix](#appendix)
   * [Git hooks](#git-hooks)

# Requirements

Go version >= 1.16

## Redis

The tasks you create are queued to Redis by a producer and those tasks are
popped from the Redis queue by a consumer. These microservices can work with
a single Redis instance or a Redis Cluster. Create either and set the credentials
in `.env`.

For a single Redis instance
```bash
REDIS_HOST=0.0.0.0
REDIS_PORT=6379
```

For Redis Cluster, set the url separated by a comma.
```bash
REDIS_ADDRESS=localhost:7001,localhost:7002,localhost:7003,localhost:7004,localhost:7005,localhost:7006
```

## (Optional) Protobuf

To use protobuf, install protobuf compiler [https://grpc.io/docs/protoc-installation/](https://grpc.io/docs/protoc-installation/) and its Go library

    go install google.golang.org/protobuf/cmd/protoc-gen-go
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

Add to `PATH`

    export PATH="$PATH:$(go env GOPATH)/bin"

## (Optional) Linting

Linting check code quality with `golangci-lint` and `gosec`

Install

    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/securego/gosec/v2/cmd/gosec@latest

# Run

 - (Optional) Database
 - Producer
 - Consumer

## Database

Optionally, set database credentials in `.env` file. It is not used in this code
base.

## Producer

We can build two kind of binaries from this same codebase, producer and consumer.

Multiple producers can be run by setting its port number to environment variable

    API_PORT=4001 go run cmd/producer/main.go
    API_PORT=4002 go run cmd/producer/main.go

## Consumer

Consumer can be run with

    go run cmd/consumer/main.go

May also set different `API_PORT` for multiple consumers.

    API_PORT=8001 go run cmd/consumer/main.go

# Gracefully Shutdown

Since both, especially consumer API contains long-running tasks, it is important
that we do not suddenly terminate those running tasks. Both producer and
consumer APIs intercepts operating system signal for termination and immediately
blocks any new requests or popping of queue from Redis. It then waits for 8
seconds to allow pending tasks to complete. If there are tasks still not
completed, it is re-enqueued to Redis to be picked by other API, or the same
API when it is restarted.

There are two ways to gracefully shutdown:

1. Press Ctrl+C in the terminal window where you started the service.


2. Send OS signal to the service PID.

Find the process ID of this API by port

    lsof -t -i:4001

Send SIGTERM signal

    kill -SIGTERM <pid number>

Or one-liner

    kill -SIGTERM $(lsof -t -i:4001)

# Usage

## Producer

Send a task by sending an http request to a Producer. In production, send the request to a load balancer that does reverse proxy between your producers.

```bash
 curl -v --location --request POST 'http://0.0.0.0:4001/api/v1/email' \
 --header 'Content-Type: application/json' \
 --data-raw '{
     "sent_by":1,
     "type":"email:deliver",
     "request_uuid": "a835b85b-6d09-4898-8bbd-1bb1406151f5", 
     "email":{
         "to":2,
         "subject":"test subject",
         "content":"test content"
     }
 }'
```

## Consumer

The consumer server(s) will automatically pick up the task from Redis.

# Monitor

Run `asynq stats` or [`asynqmon`](#asynqmon) [http://localhost:8080/](http://localhost:8080/) and, you will see tasks being process, succeeded, or failed tasks.

![inspect-payload](assets/asynqmon-inspect-payload.png)
Payload can be inspected in this asynqmon web UI if payload is serialized with JSON.

## asynqmon

Install

    git clone https://github.com/hibiken/asynqmon.git
    cd asynqmon
    make build


Ensure `/usr/local/bin` is in your path. Otherwise, add the following line to your `bashrc`

    sudo mv ./asynqmon /usr/local/bin
    source ~/.bashrc
    PATH=$PATH:/usr/local/bin

Run

    asynqmon

# Create a New Task

To create a new task, there are two ways to do it, using CLI command or manually.

To use CLI, install the program with

    # If you haven't
    git clone git@github.com:gmhafiz/asynq.git && cd asynq

    make cli

Then move the binary to your `$PATH`. For example

    mv ./bin/asynqgen ~/.local/bin/asynqgen

Check with

    asynqgen version

To create a new task, you will need to define which domain it belongs to. This 
groups related tasks within the same domain. We use `task` sub command to create
it.

For example, client requests a .csv export of all users. Domain is `export` and
task name is `users`. Thus,

    asynqgen task --domain export --name users

`asynqgen` automatically creates a new domain in `./internal/domain/export` 
along with its files. These include
 - `./internal/domain/export/handler/http/register.go`
   - registers Rest API endpoint `POST /api/v1/export`
   - register the handler this endpoint goes to 
 - `./internal/domain/export/handler/http/handler.go`
   - receives and route endpoint
   - converts Json request to local `Request` struct
   - does validation
   - depending on task type, send of to business logic layer
 - `./internal/domain/export/usecase/usecase.go`
   - does business logic
   - queues task to Redis
 - `./internal/domain/export/request.go`
   - contains domain-specific `Export` struct
   - add fields to the domain struct to define payload needed in the task
 - `./internal/task/export/create.go`
   - helper function to instantiate an `asynq` `Export` task
 - `./internal/task/export/process.go`
   - Actual *Worker* that process tasks. It is instantiated by `NewExportProcessor()`
   and processed in `ProcessTask(ctx context.Context, task *asynq.Task) error`

It then adds a new task constant in `./internal/task/tasks.go`
```go
const TypeExportUsers = "export:users"
```

`register.go` is added to `./internal/server/initDomains.go`. This adds a new 
`Export` domain to the app. Dependency injection is also done here.

Finally, it adds the task constant to `./cmd/consumer/main.go`
```go
 mux.Handle(task.TypeExportUsers, export.NewExportProcessor(s.DB()))
```
By registering the mux, consumer will be able to pop queues from Redis. This is
also the place where you inject dependencies to a consumer (worker).

To test, start both Producer and Consumer, then send a new task to the Producer 
Rest endpoint. 

```bash
go run cmd/producer/main.go
go run cmd/consumer/main.go

curl -v --location --request POST 'http://0.0.0.0:4001/api/v1/export' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "sent_by":1,
        "type":"export:users",
        "request_uuid": "a835b85b-6d09-4898-8bbd-1bb1406151f6",
        "export": {}
    }'
```

and you'll get a 200 response after 20 seconds
```bash
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Mon, 03 Jan 2022 07:20:41 GMT
< Content-Length: 0
< 
```

You may add additional task (e.g. `roles`) under the same domain with

      asynqgen task --domain export --name roles


# Appendix

## Git hooks

Add git hooks by registering the path to `.git-hooks` folder

    git config core.hooksPath .git-hooks
