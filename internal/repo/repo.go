package repo

import "github.com/ozoncp/ocp-request-api/internal/models"

// Repo is a Requests storage
type Repo interface {
	Add(requests []models.Request) error
	List(limit, offset uint64) ([]models.Request, error)
	Describe(id uint64) (*models.Request, error)
}
