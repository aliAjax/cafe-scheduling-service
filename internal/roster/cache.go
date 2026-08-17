package roster

import "sync"

type Cache struct {
	mu    sync.RWMutex
	weeks map[string]WeekView
}

func NewCache() *Cache {
	return &Cache{weeks: make(map[string]WeekView)}
}

func (c *Cache) Get(weekStart string) (WeekView, bool) {
	view, ok := c.weeks[weekStart]
	if !ok {
		return WeekView{}, false
	}
	return cloneWeek(view), true
}

func (c *Cache) Put(view WeekView) {
	c.weeks[view.WeekStart] = cloneWeek(view)
}

func (c *Cache) Delete(weekStart string) {
	delete(c.weeks, weekStart)
}
