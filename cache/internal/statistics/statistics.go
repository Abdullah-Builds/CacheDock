package statistics

import "sync/atomic"

type Stats struct {
	Hits      atomic.Uint64
	Misses    atomic.Uint64
	Requests  atomic.Uint64
	Sets      atomic.Uint64
	Deletes   atomic.Uint64
	Evictions atomic.Uint64
}
