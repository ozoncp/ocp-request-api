ifeq ($(DB_DSN),)
DB_DSN := "postgres://postgres:12345@localhost:5432/postgres?sslmode=disable"  # default for migrate command
endif

.PHONY: build
build: vendor-proto .generate .build

.PHONY: .generate
.generate:
		mkdir -p swagger
		mkdir -p pkg/ocp-request-api
		buf generate api/ocp-request-api
		mv pkg/ocp-request-api/github.com/ozoncp/ocp-request-api/pkg/ocp-request-api/* pkg/ocp-request-api/
		rm -rf pkg/ocp-request-api/github.com
		mkdir -p cmd/ocp-request-api

.PHONY: mockgen
mockgen:
		go generate internal/mockgen.go

.PHONY: .build
.build:
		CGO_ENABLED=0 GOOS=linux go build -o bin/ocp-request-api cmd/ocp-request-api/main.go

.PHONY: vendor-proto
vendor-proto: .vendor-proto

.PHONY: .vendor-proto
.vendor-proto:
		mkdir -p vendor.protogen
		mkdir -p vendor.protogen/api/ocp-request-api
		cp api/ocp-request-api/ocp-request-api.proto vendor.protogen/api/ocp-request-api
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/github.com/envoyproxy ]; then \
			mkdir -p vendor.protogen/github.com/envoyproxy &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/github.com/envoyproxy/protoc-gen-validate ;\
		fi


.PHONY: deps
deps: install-go-deps

.PHONY: install-go-deps
install-go-deps: .install-go-deps

.PHONY: .install-go-deps
.install-go-deps:
		ls go.mod || go mod init
		go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
		go get -u github.com/golang/protobuf/proto
		go get -u github.com/golang/protobuf/protoc-gen-go
		go get -u google.golang.org/grpc
		go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
		go get -u github.com/envoyproxy/protoc-gen-validate
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
		go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
		go install github.com/envoyproxy/protoc-gen-validate

.PHONY: migrate
migrate: .install-migrate-deps .migrate

.PHONY: .install-migrate-deps
.install-migrate-deps:
		go get -u github.com/pressly/goose/v3/cmd/goose

.PHONY: .migrate
.migrate:
		goose -dir sql-migrations postgres $(DB_DSN) up

.PHONY: test
test:
		go test ./... -v -coverprofile cover.out
		go tool cover -func cover.out
