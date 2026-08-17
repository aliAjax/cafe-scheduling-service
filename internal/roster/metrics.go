package roster

import "sync"

type Metrics struct {
	mu      sync.Mutex
	imports int64
	reads   int64
}

func (m *Metrics) RecordImport(count int) {
	m.imports += int64(count)
}

func (m *Metrics) RecordRead() {
	m.reads++
}

func (m *Metrics) Snapshot() (imports, reads int64) {
	return m.imports, m.reads
}
