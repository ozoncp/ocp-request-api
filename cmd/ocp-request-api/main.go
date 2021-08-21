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

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ozoncp/ocp-request-api/internal/api"
	"google.golang.org/grpc"

	desc "github.com/ozoncp/ocp-request-api/pkg/ocp-request-api"
)

const (
	grpcPort           = ":82"
	grpcServerEndpoint = "localhost:82"
)

func mustGetEnvVar(name, errorMsgIfNotDefined string) string {
	envVal := os.Getenv(name)
	if envVal == "" {
		panic(errorMsgIfNotDefined)
	}
	return envVal
}

func run(database *sql.DB) error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	desc.RegisterOcpRequestApiServer(s, api.NewRequestApi(repo.NewRepo(database)))

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
	dsn := mustGetEnvVar("OCP_REQUEST_DSN",
		"Cannot connect to db. OCP_REQUEST_DSN is not set")
	database := db.Connect(dsn)
	defer database.Close()

	go runJSON()
	if err := run(database); err != nil {
		log.Fatal(err)
	}
}
