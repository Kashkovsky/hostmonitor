package core

import "sync"

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
				Id:         prev.Id,
				InProgress: true,
				Tcp:        prev.Tcp,
				HttpStatus: prev.HttpStatus,
				Duration:   prev.Duration,
			})
		}
	} else {
		s.results.Store(res.Id, res)
	}
}

func (s *Store) Clear() {
	s.results = sync.Map{}
}

func (s *Store) ForEach(f func(TestResult) bool) {
	s.results.Range(func(key, value any) bool {
		return f(value.(TestResult))
	})
}
