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
// If number of items in the last chunk is less than Flusher's chunk size this will not be stored,
// but returned as method result.
func (f flusher) Flush(requests []models.Request) ([]models.Request, error) {
	var err error
	remains := make([]models.Request, 0, f.chunkSize)
	for ix, chunk := range utils.SplitToBulks(requests, f.chunkSize) {
		if len(chunk) == int(f.chunkSize) {
			if err = f.requestRepo.Add(chunk); err != nil {
				remains = append(remains, requests[ix*int(f.chunkSize):]...)
				return remains, err // partially added
			}
		} else { // last chunk that still should be kept in buffer
			remains = append(remains, chunk...)
		}
	}
	return remains, nil
}
