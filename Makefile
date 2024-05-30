format:
	./scripts/git-commit-format.sh

templ-gen:
	templ generate

build:
	go build -o ./tmp/main ./cmd/main

migrate-local:
	goose -dir ./internal/db/migrations turso file:./test.db up

migrate-prod:
	goose -dir ./internal/db/migrations turso ${TURSO_DATABASE_URL}?authToken=${TURSO_AUTH_TOKEN} up

seed:
	goose -dir ./internal/db/seeds -no-versioning turso file:./test.db up

build-prod:
	GOOS=linux GOARCH=amd64 go build -o aats ./cmd/main
