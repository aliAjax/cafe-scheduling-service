package roster

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Planner struct {
	store   *Store
	cache   *Cache
	metrics *Metrics

	// weekLocks serializes cache invalidation in Import against the
	// read-build-fill cycle in Week, keyed by week start. Without it, a
	// concurrent Import can delete a cache entry in the window between
	// Week's store.List and cache.Put, so Week repopulates the cache with
	// a stale snapshot that subsequent reads observe.
	weekLocks sync.Map // map[string]*sync.Mutex
}

func NewPlanner(store *Store, cache *Cache, metrics *Metrics) *Planner {
	return &Planner{store: store, cache: cache, metrics: metrics}
}

func (p *Planner) weekMutex(weekStart string) *sync.Mutex {
	v, _ := p.weekLocks.LoadOrStore(weekStart, &sync.Mutex{})
	return v.(*sync.Mutex)
}

func (p *Planner) Import(ctx context.Context, assignments []Assignment) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := p.store.Save(ctx, assignments); err != nil {
		return fmt.Errorf("save assignments: %w", err)
	}
	// Invalidate each affected week under its week lock so the delete
	// cannot interleave with a concurrent Week's read-build-fill for the
	// same week, which would let a stale snapshot back into the cache.
	seen := make(map[string]struct{}, len(assignments))
	for _, assignment := range assignments {
		date, err := time.Parse(dateLayout, assignment.WorkDate)
		if err != nil {
			return fmt.Errorf("parse imported work date: %w", err)
		}
		monday := formatDate(date.AddDate(0, 0, -((int(date.Weekday()) + 6) % 7)))
		if _, dup := seen[monday]; dup {
			continue
		}
		seen[monday] = struct{}{}
		mu := p.weekMutex(monday)
		mu.Lock()
		p.cache.Delete(monday)
		mu.Unlock()
	}
	p.metrics.RecordImport(len(assignments))
	return nil
}

func (p *Planner) Week(ctx context.Context, weekStart string) (WeekView, error) {
	if err := ctx.Err(); err != nil {
		return WeekView{}, err
	}
	mu := p.weekMutex(weekStart)
	mu.Lock()
	defer mu.Unlock()

	if view, ok := p.cache.Get(weekStart); ok {
		p.metrics.RecordRead()
		return view, nil
	}
	assignments, err := p.store.List(ctx)
	if err != nil {
		return WeekView{}, fmt.Errorf("list assignments: %w", err)
	}
	view, err := BuildWeekView(ctx, weekStart, assignments)
	if err != nil {
		return WeekView{}, fmt.Errorf("build week view: %w", err)
	}
	p.cache.Put(view)
	p.metrics.RecordRead()
	return cloneWeek(view), nil
}
