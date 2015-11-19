subdirectories = ./client ./server ./utils ./packet ./conn

all: examples build test 

build:
	go build $(subdirectories)

test:
	go test $(subdirectories)

examples: ./examples/server/server ./examples/client/client

./examples/server/server: ./examples/server/server.go
	go build -o $@ $^

./examples/client/client: ./examples/client/client.go
	go build -o $@ $^


clean:
	go clean ./...
