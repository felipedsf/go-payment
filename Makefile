### -------------- API GATEWAY ---------------- ###
tidy-gateway:
	@cd ./api_gateway && go mod tidy
build-gateway: tidy-gateway
	@go build -o ./api_gateway/bin/gateway ./api_gateway/cmd/main.go
run-gateway: build-gateway
	@./api_gateway/bin/gateway

### -------------- AUTH ---------------- ###
tidy-auth:
	@cd ./auth && go mod tidy
build-auth: tidy-auth proto-auth
	@go build -o ./auth/bin/auth ./auth/cmd/main.go
run-auth: build-auth
	@./auth/bin/auth
proto-auth:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative auth/proto/auth_svc.proto

### -------------- MONEY MOVEMENT ---------------- ###
tidy-money:
	@cd ./money_movement && go mod tidy
build-money: tidy-money proto-money
	@go build -o ./money_movement/bin/money ./money_movement/cmd/main.go
run-money: build-money
	@./money_movement/bin/money
proto-money:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative money_movement/proto/money_movement_svc.proto

### -------------- EMAIL ---------------- ###
tidy-email:
	@cd ./email && go mod tidy
build-email: tidy-email
	@go build -o ./email/bin/email ./email/cmd/consumer.go

run-email: build-email
	@./email/bin/email

### -------------- LEDGER ---------------- ###
tidy-ledger:
	@cd ./ledger && go mod tidy
build-ledger: tidy-ledger
	@go build -o ./ledger/bin/ledger ./ledger/cmd/consumer.go
run-ledger: build-ledger
	@./ledger/bin/ledger

### -------------- INFRA ---------------- ###
infra-start:
	docker-compose up -d
infra-stop:
	docker-compose down

global-tidy: tidy-gateway tidy-money tidy-auth tidy-email tidy-ledger

.PHONY: build-gateway run-gateway
