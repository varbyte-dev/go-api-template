# ==============================================================================
# go-api-template — Makefile
# ==============================================================================

APP = ./cmd/main.go
BIN = ./bin/api
GO  = /usr/local/go/bin/go

.PHONY: help run build test test-race tidy lint clean air docker-up docker-down

## help: muestra los targets disponibles
help:
	@echo "Uso: make <target>"
	@echo ""
	@grep -E '^## [a-zA-Z_-]+:' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {sub(/^## /, ""); printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

## run: ejecuta la aplicación directamente con go run
run:
	$(GO) run $(APP)

## build: compila el binario en ./bin/api (requiere gcc para CGO/SQLite)
build:
	@mkdir -p ./bin
	CGO_ENABLED=1 $(GO) build -ldflags="-s -w" -o $(BIN) $(APP)
	@echo "Binario compilado en $(BIN)"

## test: ejecuta todos los tests con cobertura
test:
	$(GO) test ./... -v -cover

## test-race: ejecuta todos los tests con detección de race conditions
test-race:
	$(GO) test ./... -race

## tidy: limpia y actualiza dependencias en go.mod / go.sum
tidy:
	$(GO) mod tidy

## lint: ejecuta golangci-lint (requiere: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	golangci-lint run ./...

## clean: elimina el binario compilado y la base de datos local
clean:
	@rm -f $(BIN) data.db
	@echo "Limpieza completada"

## air: hot-reload para desarrollo (requiere: go install github.com/air-verse/air@latest)
air:
	air

## docker-up: construye la imagen y levanta los servicios con Docker Compose
docker-up:
	docker compose up --build

## docker-down: detiene y elimina los contenedores de Docker Compose
docker-down:
	docker compose down
