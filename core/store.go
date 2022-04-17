package core

import (
	"sync"
)

type Store struct {
	results sync.Map
}

func NewStore() Store {
	syncMap := sync.Map{}

	return Store{
		results: syncMap,
	}
}

func (s *Store) AddOrUpdate(res TestResult) {
	if res.InProgress {
		existing, ok := s.results.Load(res.Id)
		if !ok {
			s.results.Store(res.Id, res)
		} else {
			prev := existing.(TestResult)
			s.results.Store(res.Id, TestResult{
				Id:           prev.Id,
				InProgress:   true,
				Tcp:          prev.Tcp,
				HttpResponse: prev.HttpResponse,
				Duration:     prev.Duration,
				Status:       prev.Status,
			})
		}
	} else {
		s.results.Store(res.Id, res)
	}
}

func (s *Store) Clear() {
	s.results.Range(func(k, _ any) bool {
		s.results.Delete(k)
		return true
	})
}

func (s *Store) ForEach(f func(TestResult) bool) {
	s.results.Range(func(key, value any) bool {
		return f(value.(TestResult))
	})
}
