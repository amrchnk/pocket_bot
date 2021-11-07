.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	docker build -t telegram-bot .

start-container:
	docker run --name tg-bot -p 8080:8080 telegram-bot
