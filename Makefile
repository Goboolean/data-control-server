build-app:
	docker compose -f ./build/docker-compose.yml up --build -d

clean-app:
	docker compose -f ./proxy-server/docker-compose.yml down

test-app:
	docker compose -f ./build/docker-compose-test.yml up redis --build -d
	if docker compose -f ./build/docker-compose-test.yml up server --build ; then \
		docker compose -f ./build/docker-compose-test.yml down; \
		exit 0; \
	else \
		docker compose -f ./build/docker-compose-test.yml down; \
		exit 1; \
	fi

benchmark-app:
	go test -bench . -cpuprofile cpu.out

sqlc-synchronize:
	curl -s -L https://raw.githubusercontent.com/Goboolean/shared/main/api/sql/schema.sql -o ./api/sql/schema.sql; \
	curl -s -L https://raw.githubusercontent.com/Goboolean/shared/main/api/sql/schema.test.sql -o ./api/sql/schema.test.sql; \

sqlc-generate: \
	sqlc-synchronize; \
	sqlc generate -f ./api/sql/sqlc.yml

sqlc-check: \
	sqlc-synchronize; \
	sqlc compile -f ./api/sql/sqlc.yml

GRPC_PROTO_PATH = ./api/grpc/fetch-server.proto
GRPC_GEN_PATH = .

grpc-generate:
	protoc \
		--go_out=${GRPC_GEN_PATH} \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GRPC_GEN_PATH) \
		--go-grpc_opt=paths=source_relative \
    ${GRPC_PROTO_PATH}	

REDIS_MODEL_PROTO_PATH = ./api/redis-model/model.proto
REDIS_MODEL_GEN_PATH = ./internal/infrastructure/redis

proto-generate:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
	./api/model/model.proto


wire-generate:
	wire cmd/inject/infrastructure.go cmd/inject/service.go

all-generate: \
	sqlc-generate \
	grpc-generate \
	proto-generate \
	wire-generate