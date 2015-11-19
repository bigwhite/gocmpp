all: examples build test 

build:
	go build ./client ./server ./utils ./packet ./conn

test:
	go test ./client ./server ./utils ./packet ./conn

examples: ./examples/server/server ./examples/client/client

./examples/server/server: ./examples/server/server.go
	go build -o $@ $^

./examples/client/client: ./examples/client/client.go
	go build -o $@ $^


clean:
	go clean ./...
