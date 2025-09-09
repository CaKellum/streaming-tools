default: build

build: 
	go build

clean: 
	rm -rf streaming-tools
	go mod tidy
	go clean

run: 
	go run main.go
