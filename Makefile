ifeq ($(POSTGRES_SETUP_MIGRATION),)
	POSTGRES_SETUP_MIGRATION := user=postgres password=postgres dbname=postgres host=localhost port=5435 sslmode=disable
endif

ifeq ($(POSTGRES_SETUP_TEST_MIGRATION),)
	POSTGRES_SETUP_TEST_MIGRATION := user=postgres password=postgres dbname=postgres host=localhost port=8888 sslmode=disable
endif


GRPC_UTILS_FOLDER:=$(CURDIR)/bin
GRPC_API_GEN=internal/gen_proto

GOOSE_PATH=$(GOPATH)/bin
MIGRATION_FOLDER=$(CURDIR)/migration/postgres

MOCKGEN_TAG=1.2.0

DOCKER_GET_KEYS_PATH=./build/gen_keys
DOCKER_GET_KEYS_IMAGE_NAME=gen_keys
DOCKER_GET_KEYS_CONTAINER_NAME=gen_keys

#GRPC
.PHONY:.install-deps
.install-deps:
	mkdir -p $(GRPC_UTILS_FOLDER)
	GOBIN=$(GRPC_UTILS_FOLDER) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(GRPC_UTILS_FOLDER) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(GRPC_UTILS_FOLDER) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2

.PHONY:.get-deps
.get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway

.PHONY:.vendor-proto
.vendor-proto:
	if [ ! -d third_party/googleapis ]; then \
		git clone https://github.com/googleapis/googleapis third_party/googleapis; \
	fi

generate-api:.install-deps  .get-deps .vendor-proto
	mkdir -p $(GRPC_API_GEN)
	protoc --proto_path api/proto --proto_path third_party/googleapis \
	--go_out=$(GRPC_API_GEN) --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=$(GRPC_UTILS_FOLDER)/protoc-gen-go \
	--go-grpc_out=$(GRPC_API_GEN) --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(GRPC_UTILS_FOLDER)/protoc-gen-go-grpc \
	--grpc-gateway_out=$(GRPC_API_GEN) --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=$(GRPC_UTILS_FOLDER)/protoc-gen-grpc-gateway \
	api/proto/*.proto

# Docker Compose
.PHONY: kafka-up
kafka-up:
	docker-compose -f deploy/docker-compose.kafka-only.yaml up --build -d

.PHONY: redis-up
redis-up:
	docker-compose -f deploy/docker-compose.cache.yaml up --build -d

.PHONY: postgresql-up
postgresql-up:
	docker-compose -f deploy/docker-compose.bd.yaml up --build -d

.PHONY: monitoring-up
monitoring-up:
	docker-compose -f deploy/docker-compose.monitoring.yaml up --build -d

.PHONY: httpservice
httpservice: kafka-up redis-up postgresql-up monitoring-up
	sleep 10
	docker-compose -f deploy/docker-compose.project.yaml up --build


# Проект: создание и запуск
.PHONY: generate-keys
generate-keys:
	docker build -t $(DOCKER_GET_KEYS_IMAGE_NAME) $(DOCKER_GET_KEYS_PATH)
	docker run --name $(DOCKER_GET_KEYS_CONTAINER_NAME) -d $(DOCKER_GET_KEYS_IMAGE_NAME) sleep infinity
	docker cp $(DOCKER_GET_KEYS_CONTAINER_NAME):/application/server.key .
	docker cp $(DOCKER_GET_KEYS_CONTAINER_NAME):/application/server.crt .


# GOOSE Migrations
.PHONY: migration-create
migration-create:
	${GOOSE_PATH}/goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: migration-up
migration-up:
	${GOOSE_PATH}/goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_MIGRATION)" up

.PHONY: migration-down
migration-down:
	${GOOSE_PATH}/goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_MIGRATION)" down

.PHONY: httpservice
httpservice:
	docker-compose -f deploy/docker-compose.project.yaml up --build

#build
.PHONY: build-project
build-project:
	go build -o main cmd/HW8/main.go

httpservice-stop:
	docker-compose -f deploy/docker-compose.kafka-only.yaml down
	docker-compose -f deploy/docker-compose.bd.yaml down
	docker-compose -f deploy/docker-compose.cache.yaml down





################### Тестирование
.PHONY: test-migration-up
test-migration-up:
	${GOOSE_PATH}/goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST_MIGRATION)" up

.PHONY: test-migration-down
test-migration-down:
	${GOOSE_PATH}/goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST_MIGRATION)" down


.PHONY: .generate-mockgen
.generate-mockgen:
	PATH="$(LOCAL_BIN):$$PATH" go generate -x -run=mockgen ./...


.PHONY: .generate-mockgen-deps
.generate-mockgen-deps:
ifeq ($(wildcard $(MOCKGEN_BIN)),)
	@GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@$(MOCKGEN_TAG)
endif

test:.generate-mockgen
	docker-compose -f deploy/docker-compose.test.yaml up --build -d
	$(info Running tests...)
	go test -v ./...
