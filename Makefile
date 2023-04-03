APP=data-control-server

DOCKERABLE=$(shell id -Gn | grep -c docker)

MAIN_PATH=cmd/main/run.go
DOCKER_COMPOSE_PATH = ./build/docker-compose.yml
SQLC_PATH = ./api/sqlc/sqlc.yml
.PROTO_PATH = ./api/grpc/sample.proto
GRPC_GEN_PATH = ./internal/infrastructure/grpc




build:

ifeq ($(DOCKERABLE), 1)
	docker run $(APP) -f $(DOCKER_COMPOSE_PATH) --build -d
else
	go build -o $(APP) $(MAIN_PATH)
	./$(APP)
endif


run:

ifeq ($(DOCKERABLE), 1)
	docker run $(APP) -f $(DOCKER_COMPOSE_PATH) --build
else
	go run $(MAIN_PATH)
endif


clean:

ifeq ($(DOCKERABLE), 1)
	docker rmi $(APP)
else
	pkill $(APP)
	rm ./$(APP)
endif


test:



sqlc-generate:
	sqlc generate -f $(SQLC_PATH)

proto-generate:
	protoc \
		--go_out=$(GRPC_GEN_PATH) \
		--go-grpc_out=$(GRPC_GEN_PATH) \
		$(.PROTO_PATH)

go-init:
	rm go.mod go.sum
	go mod init github.com/Goboolean/data-control-server
	go mod tidy