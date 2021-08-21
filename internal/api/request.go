package api

import (
	"context"
	"errors"
	"github.com/ozoncp/ocp-request-api/internal/models"
	repository "github.com/ozoncp/ocp-request-api/internal/repo"
	desc "github.com/ozoncp/ocp-request-api/pkg/ocp-request-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewRequestApi creates Request API instance
func NewRequestApi(r repository.Repo) *RequestAPI {
	return &RequestAPI{repo: r}
}

type RequestAPI struct {
	desc.UnimplementedOcpRequestApiServer
	repo repository.Repo
}

// ListRequestV1 returns a list of user Requests
func (r *RequestAPI) ListRequestV1(ctx context.Context, req *desc.ListRequestsV1Request) (*desc.ListRequestsV1Response, error) {
	log.Printf("Got list request: %v", req)

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	reqs, err := r.repo.List(ctx, req.Limit, req.Offset)

	if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	}

	ret := make([]*desc.Request, 0, len(reqs))

	for _, r := range reqs {
		ret = append(ret, &desc.Request{
			Id:     r.Id,
			UserId: r.UserId,
			Type:   r.Type,
			Text:   r.Text,
		})
	}
	return &desc.ListRequestsV1Response{
		Requests: ret,
	}, nil
}

// DescribeRequestV1 returns detailed Request information by its ID
func (r *RequestAPI) DescribeRequestV1(ctx context.Context, req *desc.DescribeRequestV1Request) (*desc.DescribeRequestV1Response, error) {
	log.Printf("Got describe request: %v", req)

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ret, err := r.repo.Describe(ctx, req.RequestId)

	if errors.Is(err, repository.NotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	}

	return &desc.DescribeRequestV1Response{
		Request: &desc.Request{
			Id:     ret.Id,
			UserId: ret.UserId,
			Type:   ret.Type,
			Text:   ret.Text,
		},
	}, nil
}

// CreateRequestV1  Creates new Request and returns its new ID
func (r *RequestAPI) CreateRequestV1(ctx context.Context, req *desc.CreateRequestV1Request) (*desc.CreateRequestV1Response, error) {
	log.Printf("Got create request: %v", req)
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	newReq := models.NewRequest(
		0,
		req.UserId,
		req.Type,
		req.Text,
	)

	if newId, err := r.repo.Add(ctx, newReq); err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	} else {
		return &desc.CreateRequestV1Response{
			RequestId: newId,
		}, nil
	}
}

// RemoveRequestV1  removes Request by its ID
func (r *RequestAPI) RemoveRequestV1(ctx context.Context, req *desc.RemoveRequestV1Request) (*desc.RemoveRequestV1Response, error) {
	log.Printf("Got remove request: %v", req)
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if found, err := r.repo.Remove(ctx, req.RequestId); err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	} else {
		return &desc.RemoveRequestV1Response{
			Found: found,
		}, nil

	}
}
