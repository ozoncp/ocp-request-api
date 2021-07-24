package utils

import (
	"fmt"
	"github.com/ozoncp/ocp-request-api/internal/models"
	"math"
)

// NonUniqueIdError error object returned if SliceToMap() receives slice of Requests with non-unique Id fields.
// DuplicateId contains Id of the first occurred duplicate Request.
type NonUniqueIdError struct {
	DuplicateId uint64
}

// Error Returns a description of the error
func (err NonUniqueIdError) Error() string {
	return fmt.Sprintf("The following id is not unique: %v", err.DuplicateId)
}

// SplitToBulks Converts a given slice of Requests to a slice of slices of Requests of a given size.
func SplitToBulks(items []models.Request, chunkSize uint) [][]models.Request {
	if chunkSize == 0 {
		return make([][]models.Request, 0)
	}

	itemsLen := uint(len(items))
	chunksNum := int(math.Ceil(float64(itemsLen) / float64(chunkSize)))

	ret := make([][]models.Request, 0, chunksNum)

	for chunkStart := uint(0); chunkStart < itemsLen; chunkStart = chunkStart + chunkSize {
		chunkEnd := chunkStart + chunkSize
		if chunkEnd > itemsLen {
			chunkEnd = itemsLen
		}

		ret = append(ret, items[chunkStart:chunkEnd])
	}
	return ret
}

// SliceToMap Converts a given slice of Requests to a map where key is Request.Id and values are corresponding requests.
// The function requires requests to have unique Id, otherwise it returns NonUniqueIdError.
func SliceToMap(items []models.Request) (map[uint64]models.Request, error) {
	byId := make(map[uint64]models.Request, len(items))

	for _, req := range items {
		if _, ok := byId[req.Id]; ok {
			return nil, NonUniqueIdError{req.Id}
		}

		byId[req.Id] = req
	}
	return byId, nil
}
