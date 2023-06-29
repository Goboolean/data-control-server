APP=fetch-server

DOCKERABLE=$(shell id -Gn | grep -c docker)

MAIN_PATH=cmd/main/run.go
DOCKER_COMPOSE_PATH = ./build/docker-compose.yml
SQLC_PATH = ./api/sqlc/sqlc.yml

GRPC_PROTO_PATH = ./api/grpc/fetch-server.proto
GRPC_GEN_PATH = ./internal/infrastructure/grpc/config

REDIS_MODEL_PROTO_PATH = ./api/redis-model/model.proto
REDIS_MODEL_GEN_PATH = ./internal/infrastructure/redis



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

grpc-generate:
	protoc \
		--go_out=${GRPC_GEN_PATH}  --go_opt=paths=source_relative \
		--go-grpc_out=$(GRPC_GEN_PATH) --go-grpc_opt=paths=source_relative \
    ${GRPC_PROTO_PATH}	
	rm ${GRPC_GEN_PATH}/fetch-server.pb.go ${GRPC_GEN_PATH}/fetch-server_grpc.pb.go
	mv ${GRPC_GEN_PATH}/api/grpc/fetch-server.pb.go ${GRPC_GEN_PATH}
	mv ${GRPC_GEN_PATH}/api/grpc/fetch-server_grpc.pb.go ${GRPC_GEN_PATH}
	rm -rf ${GRPC_GEN_PATH}/api

proto-generate:
	protoc \
		--go_out=$(REDIS_MODEL_GEN_PATH) \
		--go_opt=paths=source_relative \
		$(REDIS_MODEL_PROTO_PATH)
	rm ./internal/infrastructure/redis/model.pb.go
	mv ./internal/infrastructure/redis/api/redis-model/model.pb.go ./internal/infrastructure/redis
	rm -rf ./internal/infrastructure/redis/api