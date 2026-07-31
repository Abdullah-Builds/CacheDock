# Unit Testing

## Scope

Unit tests validate individual cache and command behavior without a network listener:

- Cache set, get, overwrite, delete, and missing-key behavior.
- TTL expiration and cleanup.
- Persistence loading and saving.
- Individual `SET`, `GET`, and `DELETE` command behavior.
- YAML configuration parsing, validation, and environment-variable overrides.

## Run

```powershell
go test ./internal/cache ./internal/commands ./config
```

Or run the complete normal suite:

```powershell
go test ./...
```

## Result

Passed on 2026-08-01.