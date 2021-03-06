package saver

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozoncp/ocp-request-api/internal/mocks"
	"github.com/ozoncp/ocp-request-api/internal/models"
	"sync"
	"time"
)

func makeRequests(num uint64) (requests []models.Request) {
	for i := uint64(0); i < num; i++ {
		requests = append(requests, models.Request{
			Id:     i,
			UserId: i,
			Type:   i,
			Text:   fmt.Sprintf("%v", i),
		})
	}
	return
}

var _ = Describe("Saver", func() {

	var (
		sav         Saver
		mockFlusher *mocks.MockFlusher
		mockCtrl    *gomock.Controller
		requests    []models.Request
		ctx         context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockCtrl = gomock.NewController(GinkgoT())
		mockFlusher = mocks.NewMockFlusher(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Saver test", func() {
		JustBeforeEach(func() {
			sav = NewSaver(10, mockFlusher, time.Second)
			requests = makeRequests(10)
		})

		It("All items are saved after Close()", func() {
			defer sav.Close()

			mockFlusher.EXPECT().
				Flush(ctx, requests).
				Return(nil, nil).
				MaxTimes(1).
				MinTimes(1)

			for _, req := range requests {
				sav.Save(req)
			}

		})

		It("Saver flushes with a one call because we wrote data quickly", func() {
			defer sav.Close()
			mockFlusher.EXPECT().
				Flush(ctx, requests).
				Return(nil, nil).
				MaxTimes(1).
				MinTimes(1)

			for _, req := range requests[:5] {
				sav.Save(req)
			}
			time.Sleep(time.Second / 2)
			for _, req := range requests[5:] {
				sav.Save(req)
			}
		})

		It("Saver flushes with a two calls because of a long pause between saves", func() {
			defer sav.Close()
			mockFlusher.EXPECT().
				Flush(ctx, requests[:5]).
				Return(nil, nil).
				MaxTimes(1).
				MinTimes(1)

			mockFlusher.EXPECT().
				Flush(ctx, requests[5:]).
				Return(nil, nil).
				MaxTimes(1).
				MinTimes(1)

			// we may pass an empty requests while sleeping
			mockFlusher.EXPECT().
				Flush(ctx, []models.Request{}).
				Return(nil, nil).
				MinTimes(0).
				MaxTimes(3)

			for _, req := range requests[:5] {
				sav.Save(req)
			}
			time.Sleep(time.Second * 2)

			for _, req := range requests[5:] {
				sav.Save(req)
			}
			time.Sleep(time.Second * 2)

		})
	})

	Context("Saver state assertions test", func() {
		JustBeforeEach(func() {
			sav = &saver{
				capacity:   10,
				flusher:    mockFlusher,
				flushQueue: make(chan models.Request, 1),
				wait:       &sync.WaitGroup{},
				flushEvery: time.Second,
			}
		})

		It("Must call init() before", func() {
			Expect(func() {
				sav.Save(models.Request{})
			}).To(PanicWith("Saver instance is not init()-ed"))

		})

		It("Cannot Save() after Close()", func() {
			mockFlusher.EXPECT().
				Flush(ctx, []models.Request{}).
				Return(nil, nil).
				MaxTimes(1)

			sav.Init()
			sav.Close()
			Expect(func() {
				sav.Save(models.Request{})
			}).To(PanicWith("Saver instance is closed"))
		})

		It("Cannot init() after Close()", func() {
			mockFlusher.EXPECT().
				Flush(ctx, []models.Request{}).
				Return(nil, nil).
				MaxTimes(1)

			sav.Init()
			sav.Close()
			Expect(func() {
				sav.Init()
			}).To(PanicWith("Saver instance is closed"))
		})

		It("Cannot Close() before init()", func() {
			Expect(func() {
				sav.Close()
			}).To(PanicWith("Saver instance is not init()-ed"))
		})

	})
})
