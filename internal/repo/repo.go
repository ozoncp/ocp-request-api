package repo

import (
	"context"
	"github.com/ozoncp/ocp-request-api/internal/models"
)

// Repo is a Requests storage
type Repo interface {
	Add(ctx context.Context, requests []models.Request) error
	List(ctx context.Context, limit, offset uint64) ([]models.Request, error)
	Describe(ctx context.Context, id uint64) (*models.Request, error)
}
