# s3bench

CLI-first S3 performance test tool (Go).

## Status

This is the initial scaffold (v1-dev):
- `validate` command (YAML config validation)
- `run` command (generates result JSON + CSV)
- `compare` command (compare two runs)

> Current `run` uses a synthetic metrics generator as placeholder.
> Next step is wiring real S3 operations (PUT/GET/DELETE) via AWS SDK v2.

## Usage

```bash
go run ./cmd/s3bench validate -c bench.yaml
go run ./cmd/s3bench run -c bench.yaml -o run1.json
go run ./cmd/s3bench run -c bench.yaml -o run2.json
go run ./cmd/s3bench compare run1.json run2.json
```
