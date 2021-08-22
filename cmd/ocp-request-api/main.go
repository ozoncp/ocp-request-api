package main

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/ozoncp/ocp-request-api/internal/db"
	"github.com/ozoncp/ocp-request-api/internal/metrics"
	prod "github.com/ozoncp/ocp-request-api/internal/producer"
	repository "github.com/ozoncp/ocp-request-api/internal/repo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ozoncp/ocp-request-api/internal/api"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jaegermetrics "github.com/uber/jaeger-lib/metrics"
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
		log.Panic().Msgf("%v is not set", name)
	}
	return envVal
}

func mustGetEnvUInt(name string) uint64 {
	val := mustGetEnvString(name)
	uintVal, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		log.Panic().Msgf("%v value must be int", name)
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
		log.Panic().Msgf("failed to connecto to Kafka brokers: %v", err)
	}

	return prod.NewProducer(kafkaTopic, producer)
}

func initTracing() {
	cfg := jaegercfg.Configuration{
		ServiceName: "ocp-request-api",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := jaegermetrics.NullFactory
	tracer, _, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)

	if err != nil {
		log.Panic().Msgf("failed to initialize jaeger: %v", err)
	}
	opentracing.SetGlobalTracer(tracer)
}

func buildRequestApi() *api.RequestAPI {
	dsn := mustGetEnvString("OCP_REQUEST_DSN")
	database := db.Connect(dsn)

	batchSize := uint(mustGetEnvUInt("OCP_REQUEST_BATCH_SIZE"))
	repo := repository.NewRepo(database)
	prom := metrics.NewMetricsReporter()
	producer := buildKafkaProducer()
	tracer := opentracing.GlobalTracer()

	return api.NewRequestApi(repo, batchSize, prom, producer, tracer)
}

func run() error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Panic().Msgf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reqApi := buildRequestApi()

	desc.RegisterOcpRequestApiServer(s, reqApi)

	if err := s.Serve(listen); err != nil {
		log.Panic().Msgf("failed to serve: %v", err)
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

func runMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":9100", nil); err != nil {
		log.Panic().Msgf("metrics endpoint failed: %v", err)
	}
}

func main() {
	initTracing()

	go runJSON()
	go runMetrics()
	if err := run(); err != nil {
		log.Panic().Msgf("service exited with error: %v", err)
	}
}
