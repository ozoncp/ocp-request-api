package utils

import (
	"github.com/ozoncp/ocp-request-api/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitToBulksRequest(t *testing.T) {
	requests := make([]models.Request, 0, 5)
	for i := uint64(0); i < uint64(cap(requests)); i++ {
		requests = append(requests, models.NewRequest(i, 100, 200, "test request"))
	}

	assert.Equal(
		t,
		[][]models.Request{{requests[0]}, {requests[1]}, {requests[2]}, {requests[3]}, {requests[4]}},
		SplitToBulks(requests, 1),
	)
	assert.Equal(t, [][]models.Request{requests}, SplitToBulks(requests, 100))
	assert.Equal(t, [][]models.Request{requests}, SplitToBulks(requests, 5))
	assert.Equal(t,
		[][]models.Request{requests[:4], {requests[4]}},
		SplitToBulks(requests, 4),
	)
	assert.Equal(t,
		[][]models.Request{requests[:3], requests[3:]},
		SplitToBulks(requests, 3),
	)
	assert.Equal(t, [][]models.Request{}, SplitToBulks(requests, 0))
}

func TestSliceOfRequestsToMap(t *testing.T) {
	req1 := models.NewRequest(1, 1, 1, "request 1")
	req2 := models.NewRequest(2, 2, 2, "request 2")
	req3 := models.NewRequest(3, 3, 3, "request 3")
	requests := []models.Request{req1, req2, req3}

	byId, err := SliceToMap(requests)

	assert.Equal(t, 3, len(byId))
	assert.Equal(t, byId[1], req1)
	assert.Equal(t, byId[2], req2)
	assert.Equal(t, byId[3], req3)
	assert.Nil(t, err)

	requests = append(requests, models.NewRequest(1, 4, 4, "request 4 (but id is 1)"))
	byId, err = SliceToMap(requests)
	assert.Equal(t, 0, len(byId))
	assert.Equal(t, uint64(1), err.(NonUniqueIdError).DuplicateId)

	byId, err = SliceToMap([]models.Request{})
	assert.Equal(t, 0, len(byId))
	assert.Nil(t, err)
}

func BenchmarkSplitToBulks(b *testing.B) {
	reqNum := 10000000
	requests := make([]models.Request, 0, reqNum)

	for i := uint64(0); i < uint64(reqNum); i++ {
		requests = append(requests, models.NewRequest(i, i, i, "test request"))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		SplitToBulks(requests, 1)
		SplitToBulks(requests, 10)
		SplitToBulks(requests, 100)
		SplitToBulks(requests, 1000)
		SplitToBulks(requests, 10000)
	}
}
