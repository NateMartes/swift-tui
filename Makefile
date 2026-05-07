.PHONY: build run clean

build:
	go build -o bin/go-swift-tui ./cmd

run: build
	./bin/go-swift-tui

clean:
	rm -rf bin
	go clean
