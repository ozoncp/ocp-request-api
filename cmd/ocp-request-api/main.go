package main

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/ozoncp/ocp-request-api/internal/api"
	"github.com/ozoncp/ocp-request-api/internal/db"
	"github.com/ozoncp/ocp-request-api/internal/metrics"
	prod "github.com/ozoncp/ocp-request-api/internal/producer"
	repository "github.com/ozoncp/ocp-request-api/internal/repo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jaegermetrics "github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	desc "github.com/ozoncp/ocp-request-api/pkg/ocp-request-api"
)

const (
	grpcPort           = ":82"
	grpcServerEndpoint = "localhost:82"
	kafkaTopic         = "ocp_request_events"
)

var serviceConfig config

type config struct {
	databaseDSN           string
	kafkaBrokers          []string
	requestWriteBatchSize uint
	jaegerHostPort        string
}

func init() {
	serviceConfig = config{
		databaseDSN:           mustGetEnvString("OCP_REQUEST_DSN"),
		kafkaBrokers:          strings.Split(mustGetEnvString("OCP_KAFKA_BROKERS"), ","),
		requestWriteBatchSize: uint(mustGetEnvUInt("OCP_REQUEST_BATCH_SIZE")),
		jaegerHostPort:        mustGetEnvString("OCP_REQUEST_JAEGER_HOST_PORT"),
	}
}

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
	brokers := serviceConfig.kafkaBrokers

	cfg := sarama.NewConfig()
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, cfg)

	if err != nil {
		log.Panic().Msgf("failed to connect to Kafka brokers: %v", err)
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
			LocalAgentHostPort: serviceConfig.jaegerHostPort,
			LogSpans:           true,
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
	database := db.Connect(serviceConfig.databaseDSN)

	batchSize := serviceConfig.requestWriteBatchSize
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
	defer reqApi.Shutdown()

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
