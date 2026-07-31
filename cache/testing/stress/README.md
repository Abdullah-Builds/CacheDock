# Stress Testing

## Scope

The opt-in stress test uses the real local TCP stack and a live listener.

- 128 simultaneous clients.
- 10,000 SET/GET rounds per client.
- 2,560,000 total application commands.
- 100-key working set to avoid intentionally skewing the result with eviction.
- Validates every response and fails on connection or protocol errors.

## Run

```powershell
$env:RUN_STRESS='1'; go test -run '^TestTCPStress$' -v ./internal/handler
```

## Verified result

On 2026-08-01, the stress test completed successfully in 11.83 seconds at 216,394 commands/s, with no response mismatches or connection failures.

This is a real TCP test on one machine. A distributed test still requires separate reachable target and load-generator hosts.