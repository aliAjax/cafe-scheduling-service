package roster

import "sync"

type Metrics struct {
	mu      sync.Mutex
	imports int64
	reads   int64
}

func (m *Metrics) RecordImport(count int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.imports += int64(count)
}

func (m *Metrics) RecordRead() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reads++
}

func (m *Metrics) Snapshot() (imports, reads int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.imports, m.reads
}
