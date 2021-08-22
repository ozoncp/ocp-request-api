package api

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/ozoncp/ocp-request-api/internal/metrics"
	"github.com/ozoncp/ocp-request-api/internal/models"
	"github.com/ozoncp/ocp-request-api/internal/producer"
	repository "github.com/ozoncp/ocp-request-api/internal/repo"
	"github.com/ozoncp/ocp-request-api/internal/utils"
	desc "github.com/ozoncp/ocp-request-api/pkg/ocp-request-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewRequestApi creates Request API instance
func NewRequestApi(r repository.Repo,
	batchSize uint,
	metricsReporter metrics.MetricsReporter,
	producer producer.Producer,
	tracer opentracing.Tracer,
) *RequestAPI {
	return &RequestAPI{
		repo:      r,
		batchSize: batchSize,
		metrics:   metricsReporter,
		producer:  producer,
		tracer:    tracer,
	}
}

type RequestAPI struct {
	desc.UnimplementedOcpRequestApiServer
	repo      repository.Repo
	batchSize uint // batch size for multi create
	metrics   metrics.MetricsReporter
	producer  producer.Producer
	tracer    opentracing.Tracer
}

// ListRequestV1 returns a list of user Requests
func (r *RequestAPI) ListRequestV1(ctx context.Context, req *desc.ListRequestsV1Request) (*desc.ListRequestsV1Response, error) {
	log.Printf("Got list request: %v", req)
	span, ctx := opentracing.StartSpanFromContext(ctx, "ListRequestV1")
	defer span.Finish()

	if err := req.Validate(); err != nil {
		r.notifyApiEvent(ctx, 0, producer.Read, err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	requests, err := r.repo.List(ctx, req.Limit, req.Offset)

	if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		r.notifyApiEvent(ctx, 0, producer.Read, err)
		return nil, err
	}

	ret := make([]*desc.Request, 0, len(requests))

	for _, req := range requests {
		ret = append(ret, &desc.Request{
			Id:     req.Id,
			UserId: req.UserId,
			Type:   req.Type,
			Text:   req.Text,
		})
		r.notifyApiEvent(ctx, req.Id, producer.Read, nil)

	}
	r.metrics.IncList(1, "ListRequestV1")
	return &desc.ListRequestsV1Response{
		Requests: ret,
	}, nil
}

// DescribeRequestV1 returns detailed Request information by its ID
func (r *RequestAPI) DescribeRequestV1(ctx context.Context, req *desc.DescribeRequestV1Request) (*desc.DescribeRequestV1Response, error) {
	log.Printf("Got describe request: %v", req)
	span, ctx := opentracing.StartSpanFromContext(ctx, "DescribeRequestV1")
	defer span.Finish()

	if err := req.Validate(); err != nil {
		r.notifyApiEvent(ctx, req.RequestId, producer.Read, err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ret, err := r.repo.Describe(ctx, req.RequestId)

	if errors.Is(err, repository.NotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	}

	r.notifyApiEvent(ctx, req.RequestId, producer.Read, err)
	r.metrics.IncRead(1, "DescribeRequestV1")

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateRequestV1")
	defer span.Finish()

	if err := req.Validate(); err != nil {
		r.notifyApiEvent(ctx, 0, producer.Create, err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	newReq := models.NewRequest(
		0,
		req.UserId,
		req.Type,
		req.Text,
	)
	newId, err := r.repo.Add(ctx, newReq)

	if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	}

	r.notifyApiEvent(ctx, newId, producer.Create, err)
	r.metrics.IncCreate(1, "CreateRequestV1")
	return &desc.CreateRequestV1Response{
		RequestId: newId,
	}, nil
}

// MultiCreateRequestV1  Creates new Request and returns its new ID
func (r *RequestAPI) MultiCreateRequestV1(ctx context.Context, req *desc.MultiCreateRequestV1Request) (*desc.MultiCreateRequestV1Response, error) {
	log.Printf("Got multi create request: %v", req)
	span, ctx := opentracing.StartSpanFromContext(ctx, "MultiCreateRequestV1")
	defer span.Finish()
	if err := req.Validate(); err != nil {
		r.notifyApiEvent(ctx, 0, producer.Create, err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	toCreate := make([]models.Request, 0, len(req.Requests))

	for _, req := range req.Requests {
		toCreate = append(toCreate, models.NewRequest(0, req.UserId, req.Type, req.Text))
	}

	newIds := make([]uint64, 0, len(req.Requests))

	for _, batch := range utils.SplitToBulks(toCreate, r.batchSize) {
		e := func() error {
			childSpan, childCtx := opentracing.StartSpanFromContext(ctx, "MultiCreateRequestV1Batch")
			defer childSpan.Finish()
			ids, err := r.repo.AddMany(ctx, batch)
			if err != nil {
				log.Error().Msgf("Failed to save requests failed with %v", err)
				r.notifyApiEvent(childCtx, 0, producer.Create, err)
				return err
			}
			newIds = append(newIds, ids...)

			r.metrics.IncCreate(uint(len(ids)), "MultiCreateRequestV1")
			for _, id := range ids {
				r.notifyApiEvent(childCtx, id, producer.Create, nil)
			}
			return nil
		}()

		if e != nil {
			return nil, e
		}
	}

	return &desc.MultiCreateRequestV1Response{
		RequestIds: newIds,
	}, nil
}

// RemoveRequestV1  removes Request by its ID
func (r *RequestAPI) RemoveRequestV1(ctx context.Context, req *desc.RemoveRequestV1Request) (*desc.RemoveRequestV1Response, error) {
	log.Printf("Got remove request: %v", req)
	span, ctx := opentracing.StartSpanFromContext(ctx, "RemoveRequestV1")
	defer span.Finish()
	if err := req.Validate(); err != nil {
		r.notifyApiEvent(ctx, req.RequestId, producer.Delete, err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := r.repo.Remove(ctx, req.RequestId)
	if errors.Is(err, repository.NotFound) {
		return nil, status.Error(codes.NotFound, "request does not exist")
	} else if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	}
	r.notifyApiEvent(ctx, req.RequestId, producer.Delete, err)
	r.metrics.IncRemove(1, "RemoveRequestV1")
	return &desc.RemoveRequestV1Response{}, nil
}

// UpdateRequestV1 updates request data
func (r *RequestAPI) UpdateRequestV1(ctx context.Context, req *desc.UpdateRequestV1Request) (*desc.UpdateRequestV1Response, error) {
	log.Printf("Got update request: %v", req)
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateRequestV1")
	defer span.Finish()
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

	r.notifyApiEvent(ctx, req.RequestId, producer.Update, err)
	r.metrics.IncUpdate(1, "UpdateRequestV1")
	return &desc.UpdateRequestV1Response{}, nil
}

func (r *RequestAPI) notifyApiEvent(ctx context.Context, requestId uint64, eventType producer.EventType, apiErr error) {
	event := producer.NewEvent(ctx, requestId, eventType, apiErr)
	r.producer.Send(event)
}
