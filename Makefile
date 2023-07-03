APP=fetch-server

MAIN_PATH=cmd/main/run.go

SQLC_PATH = ./api/sqlc/sqlc.yml

GRPC_PROTO_PATH = ./api/grpc/fetch-server.proto
GRPC_GEN_PATH = ./internal/infrastructure/grpc/config

REDIS_MODEL_PROTO_PATH = ./api/redis-model/model.proto
REDIS_MODEL_GEN_PATH = ./internal/infrastructure/redis

build-app:
	docker-compose -f ./build/docker-compose.yml up --build -d

clean-app:
	docker compose -f ./proxy-server/docker-compose.yml down

test-app:
	docker-compose -f ./build/docker-compose-test.yml up redis --build -d
	if docker-compose -f ./build/docker-compose-test.yml up server --build ; then \
		docker-compose -f ./build/docker-compose-test.yml down; \
		exit 1; \
	else \
		docker-compose -f ./build/docker-compose-test.yml down; \
		cp -r ./* ../ \
		exit 0; \
	fi

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