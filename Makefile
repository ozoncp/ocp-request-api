ifeq ($(DB_DSN),)
DB_DSN := "postgres://ocp-request:12345@localhost:5432/ocp-request?sslmode=disable"  # default for migrate command
endif

.PHONY: build
build: vendor-proto .generate .build

.PHONY: .generate
.generate:
		mkdir -p swagger
		mkdir -p pkg/ocp-request-api
		protoc -I vendor.protogen \
				--go_out=pkg/ocp-request-api --go_opt=paths=import \
				--go-grpc_out=pkg/ocp-request-api --go-grpc_opt=paths=import \
				--grpc-gateway_out=pkg/ocp-request-api \
				--grpc-gateway_opt=logtostderr=true \
				--grpc-gateway_opt=paths=import \
				--validate_out lang=go:pkg/ocp-request-api \
				--swagger_out=allow_merge=true,merge_file_name=api:swagger \
				api/ocp-request-api/ocp-request-api.proto
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
		go test internal/flusher/* -v
		go test internal/saver/* -v
		go test internal/utils/* -v
		go test internal/repo/* -v
