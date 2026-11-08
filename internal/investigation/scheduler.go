package investigation

import (
	"context"
	"log"
	"sync"
)

type Scheduler struct {
	mu             sync.Mutex // kind of future planning here with the mutex but, realistically completely extra
	investigations map[Key]Investigation
	detectionsQ    chan Detection
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		investigations: make(map[Key]Investigation),
		detectionsQ:    make(chan Detection, 64), // channel as q to pass data from goroutines size 64 Detections
	}
}

func (s *Scheduler) recordDetection(d Detection) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if investigation, exists := s.investigations[d.Key]; exists {
		if d.ObservedAt.After(investigation.LastSeen) {
			investigation.LastSeen = d.ObservedAt
			s.investigations[d.Key] = investigation
		}
		return false
	}

	s.investigations[d.Key] = Investigation{
		Key:       d.Key,
		Namespace: d.Namespace,
		PodName:   d.PodName,
		FirstSeen: d.ObservedAt,
		LastSeen:  d.ObservedAt,
	}
	return true
}

func (s *Scheduler) Submit(ctx context.Context, d Detection) error {
	select {
	case s.detectionsQ <- d:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Wait for a detection or shutdown. Keep this select without a default
// so the worker blocks when the q is empty instead of polling
func (s *Scheduler) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case d := <-s.detectionsQ:
			if s.recordDetection(d) {
				log.Printf("opened investigation for %+v", d.Key)
			}
		}
	}
}
