package roster

import (
	"context"
	"sort"
	"sync"
)

type Store struct {
	mu          sync.RWMutex
	assignments map[int64]Assignment
}

func NewStore() *Store {
	return &Store{assignments: make(map[int64]Assignment)}
}

func (s *Store) Save(ctx context.Context, assignments []Assignment) error {
	validated := make([]Assignment, len(assignments))
	for i, assignment := range assignments {
		if err := validateAssignment(ctx, assignment); err != nil {
			return err
		}
		validated[i] = cloneAssignment(assignment)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, assignment := range validated {
		s.assignments[assignment.ID] = cloneAssignment(assignment)
	}
	return nil
}

func (s *Store) List(ctx context.Context) ([]Assignment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]Assignment, 0, len(s.assignments))
	for _, assignment := range s.assignments {
		items = append(items, cloneAssignment(assignment))
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, nil
}
