include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

migrate-up:
	@migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	@migrate -path ./migrations -database "$(DB_URL)" down

service-run:
	@go run ./cmd/app/main.go

# swagger-init:
# 	@echo "Generating Swagger docs..."
# 	@swag init \
# 		-g ./cmd/app/main.go \      
# 		-d ./cmd/app,./internal \   
# 		-o ./cmd/app/docs \        
# 		--parseDependency          
# 	@echo "Swagger docs generated in ./cmd/app/docs"
