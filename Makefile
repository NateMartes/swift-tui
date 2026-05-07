.PHONY: build run clean

build:
	go build -o bin/swift-tui ./cmd

run: build
	./bin/swift-tui

clean:
	rm -rf bin
	go clean
