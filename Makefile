all: examples build test 

build:
	go build 
	go build ./utils

test:
	go test
	go test ./utils

examples: ./examples/cmpp3-server/server ./examples/cmpp2-client/client ./examples/cmpp3-client/client

./examples/cmpp3-server/server: ./examples/cmpp3-server/server.go
	go build -o $@ $^

./examples/cmpp2-client/client: ./examples/cmpp2-client/client.go
	go build -o $@ $^

./examples/cmpp3-client/client: ./examples/cmpp3-client/client.go
	go build -o $@ $^

clean:
	go clean ./...
