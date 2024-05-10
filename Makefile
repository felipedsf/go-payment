### -------------- API GATEWAY ---------------- ###
build-gateway:
	@go build -o ./api_gateway/bin/gateway ./api_gateway/cmd/main.go

run-gateway: build-gateway
	@./api_gateway/bin/gateway

### -------------- AUTH ---------------- ###

build-auth: proto-auth
	@go build -o ./auth/bin/auth ./auth/cmd/main.go

run-auth: build-auth
	@./auth/bin/auth

proto-auth:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative auth/proto/auth_svc.proto

### -------------- MONEY MOVEMENT ---------------- ###

build-money: proto-money
	@go build -o ./money_movement/bin/money ./money_movement/cmd/main.go

run-money: build-money
	@./money_movement/bin/money

proto-money:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative money_movement/proto/money_movement_svc.proto

### -------------- INFRA ---------------- ###

infra-start:
	docker-compose up -d

infra-stop:
	docker-compose down

.PHONY: build-gateway run-gateway
