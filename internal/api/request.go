package api

import (
	"context"
	"errors"
	db2 "github.com/ozoncp/ocp-request-api/internal/db"
	"github.com/ozoncp/ocp-request-api/internal/models"
	repository "github.com/ozoncp/ocp-request-api/internal/repo"
	desc "github.com/ozoncp/ocp-request-api/pkg/ocp-request-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewRequestApi creates Request API instance
func NewRequestApi() *RequestAPI {
	return &RequestAPI{}
}

type RequestAPI struct {
	desc.UnimplementedOcpRequestApiServer
}

func repoFromContext(ctx context.Context) repository.Repo {
	db := db2.FromContext(ctx)
	return repository.NewRepo(db)
}

// ListRequestV1 returns a list of user Requests
func (r *RequestAPI) ListRequestV1(ctx context.Context, req *desc.ListRequestsV1Request) (*desc.ListRequestsV1Response, error) {
	log.Printf("Got list request: %v", req)
	repo := repoFromContext(ctx)
	if err := req.Validate(); err != nil {
		return nil, err
	}
	reqs, err := repo.List(ctx, req.Limit, req.Offset)

	if err != nil {
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

// DescribeTaskV1 returns detailed Request information by its ID
func (r *RequestAPI) DescribeTaskV1(ctx context.Context, req *desc.DescribeRequestV1Request) (*desc.DescribeTaskV1Response, error) {
	log.Printf("Got describe request: %v", req)

	if err := req.Validate(); err != nil {
		return nil, err
	}
	repo := repoFromContext(ctx)

	ret, err := repo.Describe(ctx, req.RequestId)

	if errors.Is(err, repository.NotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, err
	}

	return &desc.DescribeTaskV1Response{
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
		return nil, err
	}

	newReq := models.NewRequest(
		0,
		req.UserId,
		req.Type,
		req.Text,
	)

	repo := repoFromContext(ctx)
	if newId, err := repo.Add(ctx, newReq); err != nil {
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
		return nil, err
	}
	repo := repoFromContext(ctx)
	if found, err := repo.Remove(ctx, req.RequestId); err != nil {
		return nil, err
	} else {
		return &desc.RemoveRequestV1Response{
			Found: found,
		}, nil

	}
}
