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
		requestApi *api.RequestAPI
		mockRepo   *mocks.MockRepo
		mockCtrl   *gomock.Controller
		ctx        context.Context
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(mockCtrl)
		ctx = repo.NewContext(context.Background(), mockRepo)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Add new item and return its id", func() {
		JustBeforeEach(func() {
			requestApi = api.NewRequestApi()
			ctx = repo.NewContext(context.Background(), mockRepo)
		})

		It("Add request with no error", func() {
			newRequestId := uint64(19)
			mockRepo.EXPECT().
				Add(ctx, gomock.Any()).
				Return(newRequestId, nil).
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

		It("Add() params validation", func() {
			_, err := requestApi.CreateRequestV1(
				ctx, &desc.CreateRequestV1Request{UserId: 0},
			)
			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid CreateRequestV1Request.UserId: value must be greater than 0"))

		})

		It("Remove() params validation", func() {
			_, err := requestApi.RemoveRequestV1(
				ctx, &desc.RemoveRequestV1Request{},
			)
			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid RemoveRequestV1Request.RequestId: value must be greater than 0"))

		})

		It("Describe() params validation", func() {
			_, err := requestApi.DescribeRequestV1(
				ctx, &desc.DescribeRequestV1Request{},
			)
			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid DescribeRequestV1Request.RequestId: value must be greater than 0"))

		})

		It("List() params validation", func() {
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

			mockRepo.EXPECT().
				List(ctx, limit, offset).
				Return(requests, nil).
				MaxTimes(1).
				MinTimes(1)
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

		removeTest := func(expectFound bool) {
			requestId := uint64(19)
			mockRepo.EXPECT().
				Remove(ctx, requestId).
				Return(expectFound, nil).
				MaxTimes(1).
				MinTimes(1)
			resp, err := requestApi.RemoveRequestV1(
				ctx, &desc.RemoveRequestV1Request{
					RequestId: requestId,
				},
			)

			Expect(resp).
				To(Equal(&desc.RemoveRequestV1Response{
					Found: expectFound,
				}))

			Expect(err).ToNot(HaveOccurred())
		}
		It("Remove existing request with no errors", func() {
			removeTest(true)
		})
		It("Remove non-existing request with no errors", func() {
			removeTest(false)
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
