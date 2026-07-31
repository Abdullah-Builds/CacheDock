# Testing Guide

This directory documents the verification performed for the cache service.

| Area | Document | Status |
| --- | --- | --- |
| Unit | [unit/README.md](unit/README.md) | Passed |
| Integration | [integration/README.md](integration/README.md) | Passed |
| Performance | [performance/README.md](performance/README.md) | Completed |
| Load | [load/README.md](load/README.md) | Completed locally |
| Stress | [stress/README.md](stress/README.md) | Passed |

The full measurement report is available at [../PERFORMANCE_TEST_REPORT.md](../PERFORMANCE_TEST_REPORT.md).

## Notes

- Tests run from the project root (`cache/`).
- `go test ./...` runs the normal unit and integration suite.
- The sustained stress test is opt-in and requires `RUN_STRESS=1`.
- Race detection was attempted but could not run because the local Windows C compiler does not support 64-bit mode. Install a supported 64-bit C compiler, then run `go test -race ./...`.