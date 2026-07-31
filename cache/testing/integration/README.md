# Integration Testing

## Scope

Integration testing validates multi-component protocol behavior with an in-memory connection:

1. A client sends `SET`, `GET`, and `DELETE` commands.
2. The connection handler reads and dispatches commands.
3. The cache is updated or queried.
4. The expected protocol responses are returned.

The covered sequence verifies `OK`, a retrieved value, deletion, and `NOT_FOUND` responses.

## Run

```powershell
go test ./internal/handler
```

## Result

Passed on 2026-08-01.