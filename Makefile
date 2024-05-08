format:
	./scripts/git-commit-format.sh

templ-gen:
	templ generate

build:
	go build -o ./tmp/main ./cmd/main

migrate:
	goose -dir ./internal/db/migrations sqlite ./goose.db up
