# Load Testing

## Scope

The load benchmark uses concurrent persistent client connections against the complete application protocol path:

- request buffering and parsing;
- command dispatch;
- cache access;
- response generation and reading.

Each round trip performs `SET shared value` followed by `GET shared`. The test uses `net.Pipe`, so it measures application-level concurrent load without local TCP-kernel overhead.

## Run

```powershell
go test -run '^$' -bench '^BenchmarkProtocolRoundTripParallel$' -benchtime=3s -benchmem -count=3 ./internal/handler
```

## Results

The three samples ranged from about 476K to 699K round trips/s. The mean was approximately 563K SET+GET round trips/s, or 1.13M individual commands/s.

For real local TCP traffic, use the stress test described in [../stress/README.md](../stress/README.md).