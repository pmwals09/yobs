include .env

.PHONY: format
format:
	./scripts/git-commit-format.sh

.PHONY: templ-gen
templ-gen:
	templ generate

.PHONY: build
build: templ-gen
	go build -o ./tmp/main ./cmd/main

.PHONY: migrate-local
migrate-local:
	goose -dir ./internal/db/migrations turso file:./test.db up

.PHONY: migrate-prod
migrate-prod:
	goose -dir ./internal/db/migrations turso ${TURSO_DATABASE_URL}?authToken=${TURSO_AUTH_TOKEN} up

.PHONY: seed
seed:
	goose -dir ./internal/db/seeds -no-versioning turso file:./test.db up

.PHONY: build-prod
build-prod: templ-gen
	GOOS=linux GOARCH=amd64 go build -o ./build ./cmd/main

.PHONY: deploy
deploy: build-prod
	gzip ./build
	scp ./build.gz ./.env ./scripts/init.sh ./scripts/aats.service admin@${LIGHTSAIL_IP}:/tmp/aats
	ssh admin@${LIGHTSAIL_IP} 'chmod +x /tmp/aats/init.sh && /tmp/aats/init.sh'
	
