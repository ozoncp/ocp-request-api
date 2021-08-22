package internal

//go:generate mockgen -destination=./mocks/flusher_mock.go -package=mocks github.com/ozoncp/ocp-request-api/internal/flusher Flusher
//go:generate mockgen -destination=./mocks/repo_mock.go -package=mocks github.com/ozoncp/ocp-request-api/internal/repo Repo
//go:generate mockgen -destination=./mocks/saver_mock.go -package=mocks github.com/ozoncp/ocp-request-api/internal/saver Saver
//go:generate mockgen -destination=./mocks/metrics_reporter_mock.go -package=mocks github.com/ozoncp/ocp-request-api/internal/metrics MetricsReporter
//go:generate mockgen -destination=./mocks/producer_mock.go -package=mocks github.com/ozoncp/ocp-request-api/internal/producer Producer
