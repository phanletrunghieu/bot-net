.PHONY: build run-bot

build:
	go build -o ./bin/boss ./cmd/boss
	go build -o ./bin/bot ./cmd/bot
	go build -o ./bin/server ./cmd/server

run-boss:
	go build -o ./bin/boss ./cmd/boss
	./bin/boss

run-bot:
	go build -o ./bin/bot ./cmd/bot
	./bin/bot

run-server:
	go build -o ./bin/server ./cmd/server
	./bin/server