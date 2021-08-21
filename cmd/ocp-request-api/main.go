package main

import (
	"context"
	sql "github.com/jmoiron/sqlx"
	"github.com/ozoncp/ocp-request-api/internal/db"
	"github.com/ozoncp/ocp-request-api/internal/repo"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ozoncp/ocp-request-api/internal/api"
	"google.golang.org/grpc"

	desc "github.com/ozoncp/ocp-request-api/pkg/ocp-request-api"
)

const (
	grpcPort           = ":82"
	grpcServerEndpoint = "localhost:82"
)

func mustGetEnvString(name string) string {
	envVal := os.Getenv(name)
	if envVal == "" {
		log.Panicf("%v is not set", name)
	}
	return envVal
}

func mustGetEnvUInt(name string) uint64 {
	val := mustGetEnvString(name)
	uintVal, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		log.Panicf("%v value must be int", name)
	}
	return uintVal
}

func run(database *sql.DB) error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	batchSize := mustGetEnvUInt("OCP_REQUEST_BATCH_SIZE")
	desc.RegisterOcpRequestApiServer(s, api.NewRequestApi(repo.NewRepo(database), uint(batchSize)))

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}

func runJSON() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterOcpRequestApiHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		panic(err)
	}
}

func main() {
	dsn := mustGetEnvString("OCP_REQUEST_DSN")
	database := db.Connect(dsn)
	defer database.Close()

	go runJSON()
	if err := run(database); err != nil {
		log.Fatal(err)
	}
}
