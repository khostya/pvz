COMPOSE=docker compose --env-file ./deployment/docker.env  -f ./deployment/docker-compose.yaml
COMPOSE_PVZ=${COMPOSE} -p pvz
COMPOSE_POSTGRES=${COMPOSE} -p postgres
DEFAULT_PG_URL=postgres://user:password@localhost:5432

.PHONY: up
up:
	${COMPOSE_PVZ} up -d --build

.PHONY: down
down:
	${COMPOSE_PVZ} down

.PHONY: up-postgres
up-postgres:
	${COMPOSE_POSTGRES} up -d --build postgres migrations

.PHONY: up-pvz
up-pvz:
	${COMPOSE_PVZ} up -d --build pvz

.PHONY: down-pvz
down-pvz:
	${COMPOSE_PVZ} down

.PHONY: test-unit
test-unit:
	go test ./internal/... ./pkg/...

.PHONY: up-postgres
down-postgres:
	${COMPOSE_POSTGRES} down

# proto
# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN:=$(CURDIR)/bin

PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

# Установка всех необходимых зависимостей
.PHONY: bin-deps
bin-deps:
	$(info Installing binary dependencies...)

	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(LOCAL_BIN) go install github.com/vburenin/ifacemaker@latest
	GOBIN=$(LOCAL_BIN) go install go.uber.org/mock/mockgen@latest

# Вендоринг внешних proto файлов
vendor-proto: vendor-proto-rm vendor-proto/google/protobuf vendor-proto/google/api vendor-proto/validate

vendor-proto-rm:
	rm -fdr 'vendor.proto' || true

# Устанавливаем proto описания google/protobuf
.PHONY: vendor-proto/google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.proto/protobuf &&\
	cd vendor.proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor.proto/google
	mv vendor.proto/protobuf/src/google/protobuf vendor.proto/google
	rm -rf vendor.proto/protobuf

.PHONY: vendor-proto/google/api
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.proto/googleapis && \
 	cd vendor.proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor.proto/google
	mv vendor.proto/googleapis/google/api vendor.proto/google
	rm -rf vendor.proto/googleapis

.PHONY: vendor-proto/validate
vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.proto/tmp && \
		cd vendor.proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor.proto/validate
		mv vendor.proto/tmp/validate vendor.proto/
		rm -rf vendor.proto/tmp

PROTO_PATH:=./api/v1/proto
PROTO_PATH_OUT:=./pkg/api/v1/proto
DOCS_PATH:=./docs/v1

.PHONY: all-generate-proto
all-generate-proto: bin-deps vendor-proto generate-proto

.PHONY: generate-proto
generate-proto:
	mkdir -p $(PROTO_PATH_OUT)
	mkdir -p $(DOCS_PATH)
	protoc -I api/v1/proto \
		-I vendor.proto \
		$(PROTO_PATH)/*.proto \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=$(PROTO_PATH_OUT) --go_opt=paths=source_relative\
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=$(PROTO_PATH_OUT) --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:$(PROTO_PATH_OUT)"
	find $(DOCS_PATH) -name '*.*.swagger.json' -delete

DEFAULT_PG_URL=postgres://user:password@localhost:5432/pvz?sslmode=disable

.PHONY: migration-up
migration-up:
	$(eval PG_URL?=$(DEFAULT_PG_URL))
	$(LOCAL_BIN)/goose -dir ./migrations postgres "$(PG_URL)" up

.PHONY: migration-down
migration-down:
	$(eval PG_URL?=$(DEFAULT_PG_URL))
	$(LOCAL_BIN)/goose -dir ./migrations postgres "$(PG_URL)" down

.PHONY: migration-status
migration-status:
	$(eval PG_URL?=$(DEFAULT_PG_URL))
	$(LOCAL_BIN)/goose -dir ./migrations postgres "$(PG_URL)" status


.PHONY: migration-create-sql
migration-create-sql:
	$(LOCAL_BIN)/goose -dir ./migrations create $(filter-out $@,$(MAKECMDGOALS)) sql
generate:
	LOCAL_BIN=$(LOCAL_BIN) go generate ./internal/... ./pkg/...