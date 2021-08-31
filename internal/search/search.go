package search

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	sql "github.com/jmoiron/sqlx"
	"github.com/ozoncp/ocp-request-api/internal/models"
)

// Searcher implements a full text search of Requests entities
type Searcher interface {
	Search(ctx context.Context, query string, offset, limit uint64) ([]models.Request, error)
}

// NewSearcher creates a new search. Current implementation performs full text search against PostgreSQL.
func NewSearcher(db *sql.DB) Searcher {
	stmtCache := sq.NewStmtCache(db)

	return &searcher{
		stmBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(stmtCache),
	}
}

type searcher struct {
	stmBuilder sq.StatementBuilderType
}

// Search searches for Request by a given `query`. Requests are ordered by a similarity "score"
func (s *searcher) Search(ctx context.Context, query string, limit, offset uint64) ([]models.Request, error) {
	q := s.stmBuilder.Select("id, user_id, type, text").
		From("requests").
		Where("to_tsvector('russian', text) @@ to_tsquery(?)", query).
		OrderByClause("ts_rank(to_tsvector('russian', text), to_tsquery(?)) desc", query).
		Offset(offset).
		Limit(limit)

	rows, err := q.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	requests := make([]models.Request, 0, limit)
	for rows.Next() {
		req := models.Request{}
		if err := rows.Scan(&req.Id, &req.UserId, &req.Type, &req.Text); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}
