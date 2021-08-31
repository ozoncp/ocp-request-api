package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/ozoncp/ocp-request-api/internal/api"
	"github.com/ozoncp/ocp-request-api/internal/db"
	"github.com/ozoncp/ocp-request-api/internal/metrics"
	prod "github.com/ozoncp/ocp-request-api/internal/producer"
	repository "github.com/ozoncp/ocp-request-api/internal/repo"
	"github.com/ozoncp/ocp-request-api/internal/search"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jaegermetrics "github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	desc "github.com/ozoncp/ocp-request-api/pkg/ocp-request-api"
)

const (
	grpcPort           = ":82"
	grpcServerEndpoint = "localhost:82"
	kafkaTopic         = "ocp_request_events"
)

var (
	serviceConfig config
	configPath    string
)

type config struct {
	General struct {
		tWriteBatchSize uint `mapstructure:"write_batch_size"`
	}

	Db struct {
		DSN string `mapstructure:"dsn"`
	}

	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
	} `mapstructure:"kafka"`

	Jaeger struct {
		AgentHostPort string `mapstructure:"agent_host_port"`
	} `mapstructure:"jaeger"`
}

func init() {
	flag.StringVar(&configPath, "c", "config.yaml", "A path to configuration file")
}

func readConfig(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetDefault("general.write_batch_size", 1000)
	for _, param := range []string{"jaeger.agent_host_port", "kafka.brokers", "db.dsn", "general"} {
		viper.BindEnv(param,
			fmt.Sprintf("OCP_REQUEST_%v", strings.ToUpper(strings.Replace(param, ".", "_", -1))))
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	// check for required settings
	for _, s := range []string{"jaeger.agent_host_port", "kafka.brokers", "db.dsn"} {
		if !viper.IsSet(s) {
			log.Panic().Msgf("%v setting is not set", s)
		}
	}

	if err := viper.Unmarshal(&serviceConfig); err != nil {
		log.Panic().Msgf("failed to load config: %v", err)
	}
}

func buildKafkaProducer() prod.Producer {
	brokers := serviceConfig.Kafka.Brokers

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
			LocalAgentHostPort: serviceConfig.Jaeger.AgentHostPort,
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

func run() error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Panic().Msgf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	database := db.Connect(serviceConfig.Db.DSN)
	defer database.Close()
	repo := repository.NewRepo(database)
	prom := metrics.NewMetricsReporter()
	producer := buildKafkaProducer()
	defer producer.Close()
	tracer := opentracing.GlobalTracer()
	searcher := search.NewSearcher(database)

	desc.RegisterOcpRequestApiServer(
		grpcServer, api.NewRequestApi(repo, serviceConfig.General.tWriteBatchSize, prom, producer, tracer, searcher),
	)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)

	go func() {
		<-sig
		log.Info().Msgf("Got SIGTERM. Stopping...")
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(listen); err != nil {
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
	flag.Parse()
	readConfig(configPath)
	initTracing()

	go runJSON()
	go runMetrics()
	if err := run(); err != nil {
		log.Panic().Msgf("service exited with error: %v", err)
	}
}
