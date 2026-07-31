# Performance Testing

## Scope

Microbenchmarks measure in-process cache efficiency and allocations:

- `BenchmarkSetOverwrite`: replaces entries within the configured 100-key benchmark working set.
- `BenchmarkGetHit`: cache-hit lookups and LRU promotion.
- `BenchmarkMixedParallel`: concurrent 75% GET / 25% SET workload.

## Run

```powershell
go test -run '^$' -bench '^(BenchmarkSetOverwrite|BenchmarkGetHit|BenchmarkMixedParallel)$' -benchtime=3s -benchmem -count=3 ./internal/cache
```

## Results

Three 3-second samples produced these averages:

| Benchmark | Mean latency | Approximate throughput |
| --- | ---: | ---: |
| SET overwrite | 252.3 ns/op | 3.96M ops/s |
| GET hit | 96.2 ns/op | 10.4M ops/s |
| Concurrent mixed cache load | 220.7 ns/op | 4.53M ops/s |

See [../../PERFORMANCE_TEST_REPORT.md](../../PERFORMANCE_TEST_REPORT.md) for allocations and per-trial results.