package roster

import (
	"context"
	"fmt"
	"time"
)

type Planner struct {
	store   *Store
	cache   *Cache
	metrics *Metrics
}

func NewPlanner(store *Store, cache *Cache, metrics *Metrics) *Planner {
	return &Planner{store: store, cache: cache, metrics: metrics}
}

func (p *Planner) Import(ctx context.Context, assignments []Assignment) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := p.store.Save(ctx, assignments); err != nil {
		return fmt.Errorf("save assignments: %v", err)
	}
	for _, assignment := range assignments {
		date, err := time.Parse(dateLayout, assignment.WorkDate)
		if err != nil {
			return fmt.Errorf("parse imported work date: %w", err)
		}
		monday := date.AddDate(0, 0, -((int(date.Weekday()) + 6) % 7))
		p.cache.Delete(formatDate(monday))
	}
	p.metrics.RecordImport(len(assignments))
	return nil
}

func (p *Planner) Week(ctx context.Context, weekStart string) (WeekView, error) {
	if err := ctx.Err(); err != nil {
		return WeekView{}, err
	}
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
