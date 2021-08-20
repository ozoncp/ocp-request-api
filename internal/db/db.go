package db

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	sql "github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"os"
)

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
