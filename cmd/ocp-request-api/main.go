package main

import (
	"context"
	sql "github.com/jmoiron/sqlx"
	"github.com/ozoncp/ocp-request-api/internal/db"
	"github.com/ozoncp/ocp-request-api/internal/metrics"
	prod "github.com/ozoncp/ocp-request-api/internal/producer"
	repository "github.com/ozoncp/ocp-request-api/internal/repo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ozoncp/ocp-request-api/internal/api"
	"google.golang.org/grpc"

	desc "github.com/ozoncp/ocp-request-api/pkg/ocp-request-api"
)

const (
	grpcPort           = ":82"
	grpcServerEndpoint = "localhost:82"
	kafkaTopic         = "ocp_request_events"
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

func buildKafkaProducer() prod.Producer {
	brokersRaw := mustGetEnvString("OCP_KAFKA_BROKERS")
	brokers := strings.Split(brokersRaw, ",")

	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {
		log.Panicf("failed to connecto to Kafka brokers: %v", err)
	}

	return prod.NewProducer(kafkaTopic, producer)
}

func run(database *sql.DB) error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	batchSize := uint(mustGetEnvUInt("OCP_REQUEST_BATCH_SIZE"))
	repo := repository.NewRepo(database)
	prom := metrics.NewMetricsReporter()
	producer := buildKafkaProducer()

	reqApi := api.NewRequestApi(repo, batchSize, prom, producer)

	desc.RegisterOcpRequestApiServer(s, reqApi)

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

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":9100", nil); err != nil {
			log.Fatalf("metrics endpoint failed: %v", err)
		}
	}()

	go runJSON()
	if err := run(database); err != nil {
		log.Fatal(err)
	}
}
