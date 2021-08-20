package repo

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	sql "github.com/jmoiron/sqlx"
	"github.com/ozoncp/ocp-request-api/internal/models"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

const ctxRepoKey = "repo"

var NotFound = errors.New("request does not exist")

// Repo is a Requests storage
type Repo interface {
	Add(ctx context.Context, request models.Request) (uint64, error)
	AddMany(ctx context.Context, request []models.Request) error
	List(ctx context.Context, limit, offset uint64) ([]models.Request, error)
	Describe(ctx context.Context, id uint64) (*models.Request, error)
	Remove(ctx context.Context, id uint64) (bool, error)
}

// NewRepo builds a new Repo from a given db connection
func NewRepo(db *sql.DB) Repo {
	stmtCache := sq.NewStmtCache(db)

	return &repo{
		stmBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(stmtCache),
	}
}

type repo struct {
	stmBuilder sq.StatementBuilderType
}

// Add stores a single Request and returns its ID
func (r *repo) Add(ctx context.Context, request models.Request) (uint64, error) {
	query := r.stmBuilder.Insert("requests").
		Columns("user_id", "type", "text").
		Suffix("RETURNING id").
		Values(request.UserId, request.Type, request.Text)
	newTaskId := uint64(0)

	if rows, err := query.QueryContext(ctx); err != nil {
		return newTaskId, err
	} else {
		rows.Next()
		if err := rows.Scan(&newTaskId); err != nil {
			return newTaskId, err
		}
	}
	return newTaskId, nil
}

// AddMany stores a batch of Requests with a single database query
func (r *repo) AddMany(ctx context.Context, requests []models.Request) error {
	query := r.stmBuilder.Insert("requests").
		Columns("user_id", "type", "text")

	for _, r := range requests {
		query = query.Values(r.UserId, r.Type, r.Text)
	}

	if _, err := query.ExecContext(ctx); err != nil {
		return err
	}
	return nil
}

// List returns a list of stored Requests
func (r *repo) List(ctx context.Context, limit, offset uint64) ([]models.Request, error) {
	query := r.stmBuilder.Select("id, user_id, type, text").
		From("requests").
		Offset(offset). //not the fastest approach but will keep as is in favor of simplicity (ability to remove objects makes it a bit complex)
		Limit(limit)

	rows, err := query.QueryContext(ctx)
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

// Describe returns a single Request by its ID
func (r *repo) Describe(ctx context.Context, id uint64) (*models.Request, error) {
	query := r.stmBuilder.Select("id, user_id, type, text").
		From("requests").
		Where("id = ?", id)
	row, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	req := models.Request{}
	if !row.Next() {
		return nil, NotFound
	}
	if err := row.Scan(&req.Id, &req.UserId, &req.Type, &req.Text); err != nil {
		return nil, err
	} else {
		return &req, nil
	}
}

// Remove deletes Request with a given ID from the database. Returns `false` if Request didn't exist.
func (r *repo) Remove(ctx context.Context, id uint64) (bool, error) {
	query := r.stmBuilder.Delete("requests").
		Where("id = ?", id)
	ret, err := query.ExecContext(ctx)
	if err != nil {
		return false, err
	}

	if rowsDeleted, err := ret.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsDeleted > 0, err
	}
}

// NewContext returns new child context with a Repo stored as a Value
func NewContext(ctx context.Context, r Repo) context.Context {
	return context.WithValue(ctx, ctxRepoKey, r)
}

// FromContext Extracts Repo instance from a given context.
// Will panic if context didn't contain database at `ctxRepoKey` key.
func FromContext(ctx context.Context) Repo {
	r, ok := ctx.Value(ctxRepoKey).(Repo)
	if !ok {
		log.Panic().
			Msg("Db unexpectedly is not presented in context")
	}
	return r
}

// NewInterceptorWithRepo builds a grpc interceptor instance that puts Repo instance to context
func NewInterceptorWithRepo(r Repo) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		return handler(NewContext(ctx, r), req)
	}
}
