all: examples build test 

build:
	go build 
	go build ./utils

test:
	go test
	go test ./utils

examples: ./examples/server/server ./examples/client/client

./examples/server/server: ./examples/server/server.go
	go build -o $@ $^

./examples/client/client: ./examples/client/client.go
	go build -o $@ $^


clean:
	go clean ./...
