package api_test

import (
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozoncp/ocp-request-api/internal/api"
	"github.com/ozoncp/ocp-request-api/internal/mocks"
	"github.com/ozoncp/ocp-request-api/internal/models"
	"github.com/ozoncp/ocp-request-api/internal/repo"
	desc "github.com/ozoncp/ocp-request-api/pkg/ocp-request-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Flusher", func() {

	var (
		requestApi   *api.RequestAPI
		mockRepo     *mocks.MockRepo
		mockCtrl     *gomock.Controller
		ctx          context.Context
		mockProm     *mocks.MockMetricsReporter
		mockProducer *mocks.MockProducer
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(mockCtrl)
		mockProm = mocks.NewMockMetricsReporter(mockCtrl)
		mockProducer = mocks.NewMockProducer(mockCtrl)
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Add new item and return its id", func() {
		JustBeforeEach(func() {
			requestApi = api.NewRequestApi(mockRepo, 2, mockProm, mockProducer)
			ctx = context.Background()
		})

		It("Add request with no error", func() {
			newRequestId := uint64(19)
			mockRepo.EXPECT().
				Add(ctx, gomock.Any()).
				Return(newRequestId, nil).
				MaxTimes(1).
				MinTimes(1)

			mockProm.EXPECT().
				IncCreate(uint(1), "CreateRequestV1").
				MaxTimes(1).
				MinTimes(1)

			mockProducer.EXPECT().
				Send(ctx, gomock.Any()).
				MaxTimes(1).
				MinTimes(1)

			resp, err := requestApi.CreateRequestV1(
				ctx, &desc.CreateRequestV1Request{UserId: 10, Type: 11, Text: "test"},
			)

			Expect(resp).
				To(Equal(&desc.CreateRequestV1Response{
					RequestId: newRequestId,
				}))

			Expect(err).ToNot(HaveOccurred())
		})

		It("Add many requests with no error", func() {
			requestsToCreate := []models.Request{
				{
					Id:     0,
					UserId: 10,
					Type:   100,
					Text:   "one",
				},
				{
					Id:     0,
					UserId: 20,
					Type:   200,
					Text:   "two",
				},
				{
					Id:     0,
					UserId: 30,
					Type:   300,
					Text:   "three",
				},
			}
			createRequests := make([]*desc.CreateRequestV1Request, 0)
			for _, r := range requestsToCreate {
				createRequests = append(createRequests, &desc.CreateRequestV1Request{
					UserId: r.UserId,
					Type:   r.Type,
					Text:   r.Text,
				})
			}

			mockRepo.EXPECT().
				AddMany(ctx, requestsToCreate[:2]).
				Return([]uint64{1, 2}, nil).
				MaxTimes(1).
				MinTimes(1)

			mockRepo.EXPECT().
				AddMany(ctx, requestsToCreate[2:]).
				Return([]uint64{3}, nil).
				MaxTimes(1).
				MinTimes(1)

			mockProm.EXPECT().
				IncCreate(uint(2), "MultiCreateRequestV1").
				MaxTimes(1).
				MinTimes(1)

			mockProm.EXPECT().
				IncCreate(uint(1), "MultiCreateRequestV1").
				MaxTimes(1).
				MinTimes(1)

			mockProducer.EXPECT().
				Send(ctx, gomock.Any()).
				MaxTimes(3).
				MinTimes(3)

			resp, err := requestApi.MultiCreateRequestV1(
				ctx, &desc.MultiCreateRequestV1Request{Requests: createRequests},
			)

			Expect(resp).
				To(Equal(&desc.MultiCreateRequestV1Response{
					RequestIds: []uint64{1, 2, 3},
				}))

			Expect(err).ToNot(HaveOccurred())
		})

		It("Add() params validation", func() {
			mockProducer.EXPECT().
				Send(ctx, gomock.Any()).
				MaxTimes(1).
				MinTimes(1)

			_, err := requestApi.CreateRequestV1(
				ctx, &desc.CreateRequestV1Request{UserId: 0},
			)

			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid CreateRequestV1Request.UserId: value must be greater than 0"))

		})

		It("Remove() params validation", func() {
			mockProducer.EXPECT().
				Send(ctx, gomock.Any()).
				MaxTimes(1).
				MinTimes(1)

			_, err := requestApi.RemoveRequestV1(
				ctx, &desc.RemoveRequestV1Request{},
			)

			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid RemoveRequestV1Request.RequestId: value must be greater than 0"))

		})

		It("Describe() params validation", func() {
			mockProducer.EXPECT().
				Send(ctx, gomock.Any()).
				MaxTimes(1).
				MinTimes(1)

			_, err := requestApi.DescribeRequestV1(
				ctx, &desc.DescribeRequestV1Request{},
			)

			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid DescribeRequestV1Request.RequestId: value must be greater than 0"))

		})

		It("List() params validation", func() {
			mockProducer.EXPECT().
				Send(ctx, gomock.Any()).
				MaxTimes(2).
				MinTimes(2)

			_, err := requestApi.ListRequestV1(
				ctx, &desc.ListRequestsV1Request{Limit: 0},
			)

			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid ListRequestsV1Request.Limit: value must be inside range (0, 10000]"))

			_, err = requestApi.ListRequestV1(
				ctx, &desc.ListRequestsV1Request{Limit: 1000000},
			)

			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid ListRequestsV1Request.Limit: value must be inside range (0, 10000]"))

		})

		It("List requests with no error", func() {
			offset, limit := uint64(10), uint64(100)
			requests := []models.Request{
				{1, 100, 1000, "one"},
				{2, 200, 2000, "two"},
				{3, 300, 3000, "three"},
			}
			mockProm.EXPECT().
				IncList(uint(1), "ListRequestV1").
				MaxTimes(1).
				MinTimes(1)

			mockRepo.EXPECT().
				List(ctx, limit, offset).
				Return(requests, nil).
				MaxTimes(1).
				MinTimes(1)

			mockProducer.EXPECT().
				Send(ctx, gomock.Any()).
				MaxTimes(3).
				MinTimes(3)

			resp, err := requestApi.ListRequestV1(
				ctx, &desc.ListRequestsV1Request{
					Offset: offset, Limit: limit,
				},
			)

			req := make([]*desc.Request, 0, len(requests))
			for _, r := range requests {
				req = append(req, &desc.Request{
					Id:     r.Id,
					UserId: r.UserId,
					Type:   r.Type,
					Text:   r.Text,
				})
			}

			Expect(resp).
				To(Equal(&desc.ListRequestsV1Response{
					Requests: req,
				}))

			Expect(err).ToNot(HaveOccurred())
		})

		It("Update existing request", func() {
			req := models.NewRequest(1, 10, 100, "one")
			mockRepo.EXPECT().
				Update(ctx, req).
				Return(nil).
				MaxTimes(1).
				MinTimes(1)

			mockProm.EXPECT().
				IncUpdate(uint(1), "UpdateRequestV1").
				MaxTimes(1).
				MinTimes(1)

			mockProducer.EXPECT().
				Send(ctx, gomock.Any()).
				MaxTimes(1).
				MinTimes(1)

			resp, err := requestApi.UpdateRequestV1(
				ctx, &desc.UpdateRequestV1Request{
					RequestId: req.Id,
					UserId:    req.UserId,
					Type:      req.Type,
					Text:      req.Text,
				},
			)
			Expect(resp).
				To(Equal(&desc.UpdateRequestV1Response{}))

			Expect(err).ToNot(HaveOccurred())
		})

		It("Remove non-existing request with no errors", func() {
			requestId := uint64(19)
			mockRepo.EXPECT().
				Remove(ctx, requestId).
				Return(repo.NotFound).
				MaxTimes(1).
				MinTimes(1)
			_, err := requestApi.RemoveRequestV1(
				ctx, &desc.RemoveRequestV1Request{
					RequestId: requestId,
				},
			)
			Expect(err).To(Equal(status.Error(codes.NotFound, "request does not exist")))
		})

		It("Remove existing request with no errors", func() {
			requestId := uint64(19)
			mockRepo.EXPECT().
				Remove(ctx, requestId).
				Return(nil).
				MaxTimes(1).
				MinTimes(1)

			mockProm.EXPECT().
				IncRemove(uint(1), "RemoveRequestV1").
				MaxTimes(1).
				MinTimes(1)

			mockProducer.EXPECT().
				Send(ctx, gomock.Any()).
				MaxTimes(1).
				MinTimes(1)

			resp, err := requestApi.RemoveRequestV1(
				ctx, &desc.RemoveRequestV1Request{
					RequestId: requestId,
				},
			)
			Expect(resp).
				To(Equal(&desc.RemoveRequestV1Response{}))

			Expect(err).ToNot(HaveOccurred())
		})
		It("Remove non-existing request with no errors", func() {
			requestId := uint64(19)
			mockRepo.EXPECT().
				Remove(ctx, requestId).
				Return(repo.NotFound).
				MaxTimes(1).
				MinTimes(1)
			_, err := requestApi.RemoveRequestV1(
				ctx, &desc.RemoveRequestV1Request{
					RequestId: requestId,
				},
			)
			Expect(err).To(Equal(status.Error(codes.NotFound, "request does not exist")))
		})

		It("Describe existing request", func() {
			req := models.Request{
				Id:     1,
				UserId: 10,
				Type:   1000,
				Text:   "one",
			}
			mockRepo.EXPECT().
				Describe(ctx, req.Id).
				Return(&req, nil).
				MaxTimes(1).
				MinTimes(1)

			mockProm.EXPECT().
				IncRead(uint(1), "DescribeRequestV1").
				MaxTimes(1).
				MinTimes(1)

			mockProducer.EXPECT().
				Send(ctx, gomock.Any()).
				MaxTimes(1).
				MinTimes(1)

			resp, err := requestApi.DescribeRequestV1(
				ctx, &desc.DescribeRequestV1Request{
					RequestId: req.Id,
				},
			)

			Expect(resp).
				To(Equal(&desc.DescribeRequestV1Response{
					Request: &desc.Request{
						Id:     req.Id,
						UserId: req.UserId,
						Type:   req.Type,
						Text:   req.Text,
					},
				}))

			Expect(err).ToNot(HaveOccurred())
		})

		It("Describe non-existing request", func() {
			requestId := uint64(19)
			mockRepo.EXPECT().
				Describe(ctx, requestId).
				Return(nil, repo.NotFound).
				MaxTimes(1).
				MinTimes(1)
			resp, err := requestApi.DescribeRequestV1(
				ctx, &desc.DescribeRequestV1Request{
					RequestId: requestId,
				},
			)

			Expect(resp).
				To(BeNil())

			Expect(err).To(Equal(status.Error(codes.NotFound, repo.NotFound.Error())))
		})

	})

})
