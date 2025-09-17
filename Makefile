include .env
MIGRATE=migrate -path=migration -database "$(DATABASE_URL)" -verbose

devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/daixiang0/gci@v0.11.2
	go get github.com/google/wire/cmd/wire@latest

generate:
	# go mod tidy
	go generate ./...
	# go mod tidy

fmt:
	gofumpt -l -w .;gci write ./

db-migrate-up:
	$(MIGRATE) up

db-migrate-down:
	$(MIGRATE) down

build:
	go build -o ./bin/$(APP_NAME) ./cmd/

run:build
	bin/$(APP_NAME)

wire:
	wire ./cmd/

golang-linter:
	golangci-lint run ./...


docker:
	sudo docker compose up --wait -d

test:
	go test -v ./...

