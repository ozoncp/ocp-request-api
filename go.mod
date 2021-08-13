module github.com/ozoncp/ocp-request-api

go 1.16

require (
	github.com/ozoncp/ocp-request-api/pkg/ocp-request-api v0.0.1
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.1 // indirect
	github.com/golang/glog v0.0.0-20210429001901-424d2337a529 // indirect
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/lyft/protoc-gen-star v0.5.3 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.14.0
	github.com/rs/zerolog v1.23.0 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210811021853-ddbe55d93216 // indirect
	google.golang.org/grpc v1.40.0 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace (
	github.com/ozoncp/ocp-request-api/pkg/ocp-request-api => ./pkg/ocp-request-api
)