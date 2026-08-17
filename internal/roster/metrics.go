package roster

import "sync/atomic"

type Metrics struct {
	imports atomic.Int64
	reads   atomic.Int64
}

func (m *Metrics) RecordImport(count int) {
	m.imports.Add(int64(count))
}

func (m *Metrics) RecordRead() {
	m.reads.Add(1)
}

func (m *Metrics) Snapshot() (imports, reads int64) {
	return m.imports.Load(), m.reads.Load()
}
