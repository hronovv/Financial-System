include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

migrate-up:
	@migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	@migrate -path ./migrations -database "$(DB_URL)" down

service-run: swagger-docs
	@go run ./cmd/app/main.go

service-go:
	@go run ./cmd/app/main.go

swagger-docs:
	@echo "Generating Swagger docs..."
	@go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/app/main.go -o cmd/app/docs --parseDependency --parseInternal
	@echo "Swagger docs generated in cmd/app/docs"
