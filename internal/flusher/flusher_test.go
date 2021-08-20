package flusher_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozoncp/ocp-request-api/internal/flusher"
	"github.com/ozoncp/ocp-request-api/internal/mocks"
	"github.com/ozoncp/ocp-request-api/internal/models"
)

var _ = Describe("Flusher", func() {

	var (
		fl       flusher.Flusher
		mockRepo *mocks.MockRepo
		mockCtrl *gomock.Controller
		ctx      context.Context
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(mockCtrl)
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Adding items with no errors. Will not return any remains.", func() {
		JustBeforeEach(func() {
			fl = flusher.NewFlusher(2, mockRepo)
		})

		It("Added all batches with a single call to repo.", func() {
			mockRepo.EXPECT().
				Add(ctx, gomock.Any()).
				Return(nil).
				MaxTimes(1).
				MinTimes(1)

			remains, err := fl.Flush(ctx, []models.Request{
				{1, 2, 3, ""},
				{2, 2, 3, ""},
			})

			Expect(remains).To(HaveLen(0))
			Expect(err).ToNot(HaveOccurred())
		})

		It("Added all batches with a 2 calls against repo.", func() {
			mockRepo.EXPECT().
				Add(ctx, gomock.Any()).
				Return(nil).
				MaxTimes(2).
				MinTimes(2)

			remains, err := fl.Flush(ctx, []models.Request{
				{1, 2, 3, ""},
				{2, 2, 3, ""},
				{3, 2, 3, ""},
				{4, 2, 3, ""},
			})
			Expect(remains).To(HaveLen(0))
			Expect(err).ToNot(HaveOccurred())
		})

		It("Added all batches with a 2 calls against repo. 1 items remained.", func() {
			mockRepo.EXPECT().
				Add(ctx, gomock.Any()).
				Return(nil).
				MaxTimes(3).
				MinTimes(3)

			remains, err := fl.Flush(ctx, []models.Request{
				{1, 2, 3, ""},
				{2, 2, 3, ""},
				{3, 2, 3, ""},
				{4, 2, 3, ""},
				{5, 2, 3, ""},
			})
			Expect(remains).To(HaveLen(0))
			Expect(err).ToNot(HaveOccurred())

		})

	})

	Context("Repo fails to add items", func() {
		JustBeforeEach(func() {
			fl = flusher.NewFlusher(2, mockRepo)
		})

		It("Failed to add all items", func() {
			mockRepo.EXPECT().
				Add(ctx, gomock.Any()).
				Return(errors.New("failed to add")).
				MaxTimes(1).
				MinTimes(1)

			requests := []models.Request{
				{1, 2, 3, ""},
				{2, 2, 3, ""},
			}
			remains, err := fl.Flush(ctx, requests)

			Expect(remains).To(Equal(requests))
			Expect(err).To(HaveOccurred())
		})

		It("Successfully added 2 of 3 batches (partially failed)", func() {

			successFullCall1 := mockRepo.EXPECT().
				Add(ctx, gomock.Any()).
				Return(nil)

			successFullCall2 := mockRepo.EXPECT().
				Add(ctx, gomock.Any()).
				Return(nil)

			failedCall := mockRepo.EXPECT().
				Add(ctx, gomock.Any()).
				Return(errors.New("failed to add"))

			gomock.InOrder(successFullCall1, successFullCall2, failedCall)

			requests := []models.Request{
				{1, 2, 3, ""},
				{2, 2, 3, ""},
				{3, 2, 3, ""},
				{4, 2, 3, ""},
				{5, 2, 3, ""},
				{6, 2, 3, ""},
				{7, 2, 3, ""},
			}
			remains, err := fl.Flush(ctx, requests)

			Expect(remains).To(Equal([]models.Request{
				{5, 2, 3, ""},
				{6, 2, 3, ""},
				{7, 2, 3, ""},
			}), "These are failed to add to repo")
			Expect(err).To(HaveOccurred())
		})

	})

})
