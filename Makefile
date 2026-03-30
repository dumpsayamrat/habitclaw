.PHONY: test build run clean

test:
	go test ./...

build:
	go build -o habitclaw .

run:
	go run .

clean:
	rm -f habitclaw
