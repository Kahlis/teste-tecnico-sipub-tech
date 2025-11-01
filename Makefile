protoc:
	protoc --go_out=./apigateway/core --go-grpc_out=./apigateway/core ./proto/movies.proto
	protoc --go_out=./movies/core --go-grpc_out=./movies/core ./proto/movies.proto
build:
	cp .env.example .env
	docker compose up --build -d
up:
	docker compose up -d
down:
	docker compose down -v --remove-orphans
e2e:
	@echo "Running end-to-end tests..."
	go test -C movies ./core/integration/test/e2e
mock:
	go test -C movies ./core/integration/test/mock/movies_test.go
clean:
	docker volume rm teste-tecnico-sipub-tech_mongodb_data
deps:
	@echo "Downloading Go dependencies..."
	go mod -C movies download
	go mod -C apigateway download
	@echo "Verifying Go dependencies..."
	go mod -C movies verify
	go mod -C apigateway verify
	@echo "Linking Go dependencies..."
	go mod -C movies tidy
	go mod -C apigateway tidy
