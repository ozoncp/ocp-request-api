package api

import (
	"context"
	desc "github.com/ozoncp/ocp-request-api/pkg/ocp-request-api"
	"github.com/rs/zerolog/log"
	"math/rand"
)

// NewRequestApi creates Request API instance
func NewRequestApi() *RequestAPI {
	return &RequestAPI{}
}

type RequestAPI struct {
	desc.UnimplementedOcpRequestApiServer
}

// ListRequestV1 returns a list of user Requests
func (r *RequestAPI) ListRequestV1(ctx context.Context, req *desc.ListRequestsV1Request) (*desc.ListRequestsV1Response, error) {
	log.Printf("Got list request: %v", req)
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return &desc.ListRequestsV1Response{
		Requests: make([]*desc.Request, 0),
	}, nil
}

// DescribeTaskV1 returns detailed Request information by its ID
func (r *RequestAPI) DescribeTaskV1(ctx context.Context, req *desc.DescribeRequestV1Request) (*desc.DescribeTaskV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	log.Printf("Got describe request: %v", req)
	return &desc.DescribeTaskV1Response{
		Request: &desc.Request{
			Id:     1,
			UserId: 1,
			Type:   1,
			Text:   "test",
		},
	}, nil
}

// CreateRequestV1  Creates new Request and returns its new ID
func (r *RequestAPI) CreateRequestV1(ctx context.Context, req *desc.CreateRequestV1Request) (*desc.CreateRequestV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	log.Printf("Got create request: %v", req)
	return &desc.CreateRequestV1Response{
		RequestId: rand.Uint64(),
	}, nil
}

// RemoveRequestV1  removes Request by its ID
func (r *RequestAPI) RemoveRequestV1(ctx context.Context, req *desc.RemoveRequestV1Request) (*desc.RemoveRequestV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	log.Printf("Got remove request: %v", req)
	return &desc.RemoveRequestV1Response{
		Found: false,
	}, nil
}
