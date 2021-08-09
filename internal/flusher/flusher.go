package flusher

import (
	"github.com/ozoncp/ocp-request-api/internal/models"
	"github.com/ozoncp/ocp-request-api/internal/repo"
	"github.com/ozoncp/ocp-request-api/internal/utils"
)

// Flusher is interface to store Requests items to a storage
type Flusher interface {
	Flush(entities []models.Request) ([]models.Request, error)
}

// NewFlusher creates a new Flusher instance that writes Requests to storage by batches of a given size
func NewFlusher(
	chunkSize uint,
	requestRepo repo.Repo,
) Flusher {
	return &flusher{
		chunkSize:   chunkSize,
		requestRepo: requestRepo,
	}
}

type flusher struct {
	chunkSize   uint
	requestRepo repo.Repo
}

// Flush stores a slice of Requests to the underlying repository. It makes requests by chunks of a certain size.
// It's returns a slice of Requests that it's failed to write.
func (f *flusher) Flush(requests []models.Request) ([]models.Request, error) {
	var err error
	if len(requests) == 0 {
		return requests, nil
	}

	remains := make([]models.Request, 0, f.chunkSize)
	for ix, chunk := range utils.SplitToBulks(requests, f.chunkSize) {
		if err = f.requestRepo.Add(chunk); err != nil {
			remains = append(remains, requests[ix*int(f.chunkSize):]...)
			return remains, err // partially added
		}
	}
	return remains, nil
}
