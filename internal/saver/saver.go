package saver

import (
	"context"
	"github.com/ozoncp/ocp-request-api/internal/flusher"
	"github.com/ozoncp/ocp-request-api/internal/models"
	"log"
	"sync"
	"time"
)

const (
	inited = 0b01
	closed = 0b10
)

// Saver instance saves Request into underlying storage.
// User must call Init() before using an instance.
// And Close() to ensure all pending item are stored.
// Closed instance cannot be used.
type Saver interface {
	Save(entity models.Request)
	Init()
	Close()
}

// NewSaver creates a new Saver instance.
// It asynchronously collects save Requests into internally slice with given `capacity`.
// It flushes Requests into underlying `flusher` with `flushEvery` periodicity.
func NewSaver(capacity uint, flusher flusher.Flusher, flushEvery time.Duration) Saver {
	s := &saver{
		capacity:   capacity,
		flusher:    flusher,
		flushQueue: make(chan models.Request, capacity),
		wait:       &sync.WaitGroup{},
		flushEvery: flushEvery,
	}
	s.Init()
	return s
}

type saver struct {
	capacity   uint
	flusher    flusher.Flusher
	flushQueue chan models.Request
	wait       *sync.WaitGroup
	state      int8 // to check if it's closed or inited
	flushEvery time.Duration
}

// Save Request into underlying storage
func (s *saver) Save(request models.Request) {
	s.mustNotBeClosed()
	s.mustBeInitialized()

	s.flushQueue <- request
}

// Init initiates saver so it's ready to Save Requests
func (s *saver) Init() {
	s.mustNotBeClosed()
	if s.isInited() {
		return
	}

	ticker := time.NewTicker(s.flushEvery)
	s.wait.Add(1)

	go func() {
		requests := make([]models.Request, 0, s.capacity)
		defer s.wait.Done()
		for {
			select {
			case req, ok := <-s.flushQueue:
				if !ok {
					s.flush(requests)
					requests = requests[:0]
					return
				} else {
					requests = append(requests, req)
				}
			case <-ticker.C:
				s.flush(requests)
				requests = requests[:0]
			}
		}

	}()
	s.state |= inited
}

// Close ensures all Requests object are processed. Closed Saver does not accept new saves.
func (s *saver) Close() {
	s.mustBeInitialized()
	if s.isClosed() {
		return
	}

	s.state |= closed
	close(s.flushQueue)
	s.wait.Wait()
}

// flushes a slice of Requests
func (s *saver) flush(requests []models.Request) {
	ctx := context.Background()
	if failedToFlushReq, err := s.flusher.Flush(ctx, requests); err != nil {
		log.Printf("failed to save %v requests: %v", len(failedToFlushReq), err)
	}
}

func (s *saver) mustNotBeClosed() {
	if s.isClosed() {
		panic("Saver instance is closed")
	}
}

func (s *saver) mustBeInitialized() {
	if !s.isInited() {
		panic("Saver instance is not Init()-ed")
	}
}

func (s *saver) isClosed() bool {
	return s.state&closed == closed
}

func (s *saver) isInited() bool {
	return s.state&inited == inited
}
