package notify

import "sync"

type Streak struct {
	mu    sync.Mutex
	n     int
	limit int
}

func NewStreak(limit int) *Streak {
	return &Streak{limit: limit}
}

func (s *Streak) OK() {
	s.mu.Lock()
	s.n = 0
	s.mu.Unlock()
}

func (s *Streak) Fail() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.n++
	if s.limit <= 0 {
		return false
	}
	return s.n%s.limit == 0
}
