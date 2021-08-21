package api

import (
	"context"
	"errors"
	"github.com/ozoncp/ocp-request-api/internal/models"
	repository "github.com/ozoncp/ocp-request-api/internal/repo"
	"github.com/ozoncp/ocp-request-api/internal/utils"
	desc "github.com/ozoncp/ocp-request-api/pkg/ocp-request-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewRequestApi creates Request API instance
func NewRequestApi(r repository.Repo, batchSize uint) *RequestAPI {
	return &RequestAPI{repo: r, batchSize: batchSize}
}

type RequestAPI struct {
	desc.UnimplementedOcpRequestApiServer
	repo      repository.Repo
	batchSize uint // batch size for multi create
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

// MultiCreateRequestV1  Creates new Request and returns its new ID
func (r *RequestAPI) MultiCreateRequestV1(ctx context.Context, req *desc.MultiCreateRequestV1Request) (*desc.MultiCreateRequestV1Response, error) {
	log.Printf("Got multi create request: %v", req)
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	toCreate := make([]models.Request, 0, len(req.Requests))

	for _, req := range req.Requests {
		toCreate = append(toCreate, models.NewRequest(0, req.UserId, req.Type, req.Text))
	}

	newIds := make([]uint64, 0, len(req.Requests))

	for _, batch := range utils.SplitToBulks(toCreate, r.batchSize) {
		ids, err := r.repo.AddMany(ctx, batch)
		if err != nil {
			log.Error().Msgf("Failed to save requests failed with %v", err)
			return nil, err
		}
		newIds = append(newIds, ids...)
	}
	return &desc.MultiCreateRequestV1Response{
		RequestIds: newIds,
	}, nil
}

// RemoveRequestV1  removes Request by its ID
func (r *RequestAPI) RemoveRequestV1(ctx context.Context, req *desc.RemoveRequestV1Request) (*desc.RemoveRequestV1Response, error) {
	log.Printf("Got remove request: %v", req)
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := r.repo.Remove(ctx, req.RequestId)
	if errors.Is(err, repository.NotFound) {
		return nil, status.Error(codes.NotFound, "request does not exist")
	} else if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	}

	return &desc.RemoveRequestV1Response{}, nil
}

// UpdateRequestV1 updates request data
func (r *RequestAPI) UpdateRequestV1(ctx context.Context, req *desc.UpdateRequestV1Request) (*desc.UpdateRequestV1Response, error) {
	log.Printf("Got update request: %v", req)
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := r.repo.Update(
		ctx, models.NewRequest(req.RequestId, req.UserId, req.Type, req.Text),
	)
	if errors.Is(err, repository.NotFound) {
		return nil, status.Error(codes.NotFound, "request does not exist")
	} else if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	}

	return &desc.UpdateRequestV1Response{}, nil
}
