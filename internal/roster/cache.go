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
	c.mu.RLock()
	defer c.mu.RUnlock()
	view, ok := c.weeks[weekStart]
	if !ok {
		return WeekView{}, false
	}
	return cloneWeek(view), true
}

func (c *Cache) Put(view WeekView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.weeks[view.WeekStart] = cloneWeek(view)
}

func (c *Cache) Delete(weekStart string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.weeks, weekStart)
}
