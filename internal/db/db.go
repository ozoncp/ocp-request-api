package db

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	sql "github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
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
