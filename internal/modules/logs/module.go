package logs

import (
	"sync"
	"time"
)

type Entry struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Module  string    `json:"module"`
	Message string    `json:"message"`
}

type Store struct {
	mu      sync.Mutex
	max     int
	entries []Entry
}

func NewStore(max int) *Store {
	if max <= 0 {
		max = 512
	}
	return &Store{max: max}
}

func (s *Store) Add(level string, module string, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, Entry{Time: time.Now(), Level: level, Module: module, Message: message})
	if len(s.entries) > s.max {
		s.entries = append([]Entry(nil), s.entries[len(s.entries)-s.max:]...)
	}
}

func (s *Store) Query(limit int) []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 || limit > len(s.entries) {
		limit = len(s.entries)
	}
	start := len(s.entries) - limit
	return append([]Entry(nil), s.entries[start:]...)
}
