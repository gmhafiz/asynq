proto:
	protoc api/v1/*.proto \
	--go_out=. \
	--go_opt=paths=source_relative \
	--proto_path=.

	protoc api/v1/*.proto \
    --proto_path=. \
	--php_out=api/v1

test:
	go test -v ./...
	go test -race ./...

lint:
	go vet ./...
	go fmt ./...
	golangci-lint run
	gosec ./...

build: cli
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./bin/producer ./cmd/producer/main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./bin/consumer ./cmd/consumer/main.go

cli:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./bin/asynqgen ./cmd/asynqgen/main.go

check: proto lint test

all: check build
	cp .env.prod ./bin/.env
