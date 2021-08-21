package repo

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	sql "github.com/jmoiron/sqlx"
	"github.com/ozoncp/ocp-request-api/internal/models"
)

var NotFound = errors.New("request does not exist")

// Repo is a Requests storage
type Repo interface {
	Add(ctx context.Context, request models.Request) (uint64, error)
	AddMany(ctx context.Context, request []models.Request) ([]uint64, error)
	List(ctx context.Context, limit, offset uint64) ([]models.Request, error)
	Describe(ctx context.Context, id uint64) (*models.Request, error)
	Remove(ctx context.Context, id uint64) error
	Update(ctx context.Context, id models.Request) error
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

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return newTaskId, err
	}
	rows.Next()
	if err := rows.Scan(&newTaskId); err != nil {
		return newTaskId, err
	}

	return newTaskId, nil
}

// AddMany stores a batch of Requests with a single database query
func (r *repo) AddMany(ctx context.Context, requests []models.Request) ([]uint64, error) {
	query := r.stmBuilder.Insert("requests").
		Columns("user_id", "type", "text").
		Suffix("RETURNING id")

	for _, r := range requests {
		query = query.Values(r.UserId, r.Type, r.Text)
	}
	rows, err := query.QueryContext(ctx)

	if err != nil {
		return nil, err
	}
	newIds := make([]uint64, 0, len(requests))
	for rows.Next() {
		id := uint64(0)
		rows.Scan(&id)
		newIds = append(newIds, id)
	}
	return newIds, nil
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

// Remove deletes Request with a given ID from the database. Returns NotFound if Request doesn't exist.
func (r *repo) Remove(ctx context.Context, id uint64) error {
	query := r.stmBuilder.Delete("requests").
		Where("id = ?", id)
	ret, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	rowsDeleted, err := ret.RowsAffected()
	if err != nil {
		return err
	} else if rowsDeleted == 0 {
		return NotFound
	}

	return nil

}

// Update updates existing request.Returns NotFound error if request doesn't exist,
func (r *repo) Update(ctx context.Context, request models.Request) error {
	query := r.stmBuilder.
		Update("requests").
		Set("user_id", request.UserId).
		Set("type", request.Type).
		Set("text", request.Text).
		Where("id = ?", request.Id)
	ret, err := query.ExecContext(ctx)

	if err != nil {
		return err
	}

	rowsUpdated, err := ret.RowsAffected()
	if err != nil {
		return err
	} else if rowsUpdated == 0 {
		return NotFound
	}
	return nil
}
