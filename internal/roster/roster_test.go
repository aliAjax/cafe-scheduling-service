package roster

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
)

func sampleAssignment(id int64) Assignment {
	note := "front counter"
	return Assignment{
		ID: id, EmployeeID: id, WorkDate: "2026-08-17",
		StartTime: "09:00", EndTime: "17:00", Note: &note,
		Tags: []string{"barista", "opening"},
	}
}

func TestWeekViewIsolation(t *testing.T) {
	input := []Assignment{sampleAssignment(2), sampleAssignment(1)}
	view, err := BuildWeekView(context.Background(), "2026-08-17", input)
	if err != nil {
		t.Fatal(err)
	}
	if len(view.Days) != 7 || len(view.Days[0].Assignments) != 2 {
		t.Fatalf("unexpected week shape: %#v", view)
	}
	if view.Days[0].Assignments[0].ID != 1 || view.Hours != 16 {
		t.Fatalf("week ordering or hours are wrong: %#v", view)
	}

	input[0].Tags[0] = "changed"
	*input[0].Note = "changed"
	if got := view.Days[0].Assignments[1]; got.Tags[0] != "barista" || *got.Note != "front counter" {
		t.Fatalf("view aliases input: %#v", got)
	}

	cache := NewCache()
	cache.Put(view)
	view.Days[0].Assignments[0].Tags[0] = "mutated"
	cached, ok := cache.Get("2026-08-17")
	if !ok || cached.Days[0].Assignments[0].Tags[0] != "barista" {
		t.Fatalf("cache aliases caller data: %#v", cached)
	}
	cached.Days[0].Assignments[0].Tags[0] = "again"
	again, _ := cache.Get("2026-08-17")
	if again.Days[0].Assignments[0].Tags[0] != "barista" {
		t.Fatalf("cache getter exposed internal data: %#v", again)
	}
}

func TestValidationErrorChain(t *testing.T) {
	bad := sampleAssignment(1)
	bad.EmployeeID = 0
	_, err := BuildWeekView(context.Background(), "2026-08-17", []Assignment{bad})
	if !errors.Is(err, ErrInvalidSchedule) {
		t.Fatalf("errors.Is failed for %v", err)
	}
	var fieldErr *FieldError
	if !errors.As(err, &fieldErr) || fieldErr.Field != "employeeId" {
		t.Fatalf("errors.As lost field details: %#v, %v", fieldErr, err)
	}

	planner := NewPlanner(NewStore(), NewCache(), &Metrics{})
	err = planner.Import(context.Background(), []Assignment{bad})
	if !errors.Is(err, ErrInvalidSchedule) || !errors.As(err, &fieldErr) {
		t.Fatalf("service lost validation error chain: %T %v", err, err)
	}
}

func TestConcurrentPlanner(t *testing.T) {
	planner := NewPlanner(NewStore(), NewCache(), &Metrics{})
	const workers = 32
	var wg sync.WaitGroup
	errCh := make(chan error, workers)
	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			assignment := sampleAssignment(id)
			if err := planner.Import(context.Background(), []Assignment{assignment}); err != nil {
				errCh <- err
				return
			}
			if _, err := planner.Week(context.Background(), "2026-08-17"); err != nil {
				errCh <- err
			}
		}(int64(i))
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatal(err)
	}
	view, err := planner.Week(context.Background(), "2026-08-17")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(view.Days[0].Assignments); got != workers {
		t.Fatalf("got %d assignments, want %d", got, workers)
	}
	imports, reads := planner.metrics.Snapshot()
	if imports != workers || reads < workers {
		t.Fatalf("metrics lost updates: imports=%d reads=%d", imports, reads)
	}
}

func TestCancellationStopsRosterWork(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	planner := NewPlanner(NewStore(), NewCache(), &Metrics{})
	if err := planner.Import(ctx, []Assignment{sampleAssignment(1)}); !errors.Is(err, context.Canceled) {
		t.Fatalf("import error = %v, want context canceled", err)
	}
	if _, err := planner.Week(ctx, "2026-08-17"); !errors.Is(err, context.Canceled) {
		t.Fatalf("week error = %v, want context canceled", err)
	}
	if _, err := BuildWeekView(ctx, "2026-08-17", []Assignment{sampleAssignment(1)}); !errors.Is(err, context.Canceled) {
		t.Fatalf("build error = %v, want context canceled", err)
	}
	var output bytes.Buffer
	if err := EncodeWeek(ctx, &output, WeekView{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("encode error = %v, want context canceled", err)
	}
	if output.Len() != 0 {
		t.Fatalf("canceled encoding wrote %d bytes", output.Len())
	}
	imports, reads := planner.metrics.Snapshot()
	if imports != 0 || reads != 0 {
		t.Fatalf("canceled work changed metrics: %d %d", imports, reads)
	}
}

func TestOptionalFieldsRemainNilSafe(t *testing.T) {
	assignment := sampleAssignment(1)
	assignment.Note = nil
	assignment.Tags = nil
	t.Run("clone", func(t *testing.T) {
		cloned := cloneAssignment(assignment)
		if cloned.Note != nil || cloned.Tags != nil {
			t.Fatalf("nil optionals changed: %#v", cloned)
		}
	})
	t.Run("week", func(t *testing.T) {
		view, err := BuildWeekView(context.Background(), "2026-08-17", []Assignment{assignment})
		if err != nil || view.Days[0].Assignments[0].Note != nil {
			t.Fatalf("week failed for nil optionals: %#v %v", view, err)
		}
	})
	t.Run("cache", func(t *testing.T) {
		cache := NewCache()
		view, _ := BuildWeekView(context.Background(), "2026-08-17", []Assignment{assignment})
		cache.Put(view)
		got, _ := cache.Get("2026-08-17")
		if got.Days[0].Assignments[0].Note != nil {
			t.Fatalf("cache changed nil note: %#v", got)
		}
	})
	t.Run("encode", func(t *testing.T) {
		view, _ := BuildWeekView(context.Background(), "2026-08-17", []Assignment{assignment})
		var output bytes.Buffer
		if err := EncodeWeek(context.Background(), &output, view); err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(output.Bytes(), []byte(fmt.Sprintf(`"employeeId":%d`, assignment.EmployeeID))) {
			t.Fatalf("unexpected JSON: %s", output.String())
		}
	})
}
