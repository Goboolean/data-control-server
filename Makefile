APP=fetch-server

MAIN_PATH=cmd/main/run.go

SQLC_PATH = ./api/sqlc/sqlc.yml

GRPC_PROTO_PATH = ./api/grpc/fetch-server.proto
GRPC_GEN_PATH = ./internal/infrastructure/grpc/config

REDIS_MODEL_PROTO_PATH = ./api/redis-model/model.proto
REDIS_MODEL_GEN_PATH = ./internal/infrastructure/rediscache

build-app:
	docker-compose -f ./build/docker-compose.yml up --build -d

clean-app:
	docker compose -f ./proxy-server/docker-compose.yml down

test-app:
	if docker-compose -f ./build/docker-compose-test.yml up --build ; then \
		docker-compose -f ./build/docker-compose-test.yml down ; \
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
	rm ./internal/infrastructure/rediscache/model.pb.go
	mv ./internal/infrastructure/rediscache/api/redis-model/model.pb.go ./internal/infrastructure/rediscache
	rm -rf ./internal/infrastructure/rediscache/api