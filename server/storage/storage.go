package storage

import (
	"fmt"
	"github.com/amankumarsinghy77/telemon/constants"
	"log"
	"sync"
	"time"
)

type InMemoryStorage struct {
	metrics map[string][]*constants.SystemMetrics
	mu      sync.Mutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		metrics: make(map[string][]*constants.SystemMetrics),
	}
}

func (s *InMemoryStorage) Store(metrics *constants.SystemMetrics) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hourAgo := time.Now().Add(-time.Hour).UnixNano()
	var recentMetrics []*constants.SystemMetrics
	for _, m := range s.metrics[metrics.Hostname] {
		if m.Timestamp > hourAgo {
			recentMetrics = append(recentMetrics, m)
		}
	}
	recentMetrics = append(recentMetrics, metrics)
	log.Println(metrics)
	s.metrics[metrics.Hostname] = recentMetrics
}

func (s *InMemoryStorage) Query(hostname string, to, from time.Time) ([]*constants.SystemMetrics, error) {
	metrics, ok := s.metrics[hostname]
	if !ok {
		return nil, fmt.Errorf("no metrics found with hostname : %s", hostname)
	}
	var res []*constants.SystemMetrics
	for _, m := range metrics {
		if m.Timestamp >= from.UnixNano() && m.Timestamp <= to.UnixNano() {
			res = append(res, m)
		}
	}
	return res, nil
}
