package part04

import (
	"log"
	"sync"
	"testing"
)

type safeNum struct {
	count int
	// For reading more and writing less,
	// it is recommended to use a read-write mutex
	// mu sync.RWMutex
	mu sync.Mutex
}

// Get get value
func (s *safeNum) Get() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.count
}

// Set set value
func (s *safeNum) Set(i int) {
	s.mu.Lock()
	s.count = i
	s.mu.Unlock()
}

func TestSafeNum(t *testing.T) {
	i := &safeNum{count: 1}
	done := make(chan struct{}, 1)
	go func() {
		i.Set(10)
		close(done)
	}()
	<-done

	log.Println("count: ", i.Get())
}
