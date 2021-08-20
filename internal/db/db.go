package db

import (
	"context"
	_ "github.com/jackc/pgx/v4/stdlib"
	sql "github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"os"
)

const ctxDbKey = "main_db"

// NewContext returns new child context with a db stored as a Value
func NewContext(ctx context.Context, db *sql.DB) context.Context {
	ctxDB := context.WithValue(ctx, ctxDbKey, db)

	return ctxDB
}

// FromContext Extracts database connection instance from a given context.
// Will panic if context didn't contain database at `ctxDbKey` key.
func FromContext(ctx context.Context) *sql.DB {
	client, ok := ctx.Value(ctxDbKey).(*sql.DB)
	if !ok {
		log.Panic().
			Msg("Db unexpectedly is not presented in context")
	}
	return client
}

// Connect connects using given connection string and returns database instance.
// Will panic on connect failure.
func Connect(DSN string) *sql.DB {
	db, err := sql.Connect("pgx", DSN)
	if err != nil {
		log.Panic().
			AnErr("error", err).
			Msg("failed to connect to db")
	}
	return db
}

// GetDSNFromENV extracts database connection string from ENV variable
func GetDSNFromENV() string {
	dsn := os.Getenv("OCP_REQUEST_DSN")

	if dsn == "" {
		log.Panic().
			Msg("Database dsn is undefined. Set OCP_REQUEST_DSN env variable.")
	}
	return dsn
}

// NewInterceptorWithDB builds a grpc interceptor instance that puts database connection instance to context
func NewInterceptorWithDB(db *sql.DB) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		return handler(NewContext(ctx, db), req)
	}
}
