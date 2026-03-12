# s3bench

CLI-first S3 performance test tool (Go).

## Status

This is the initial scaffold (v1-dev):
- `validate` command (YAML config validation)
- `run` command (generates result JSON + CSV)
- `compare` command (compare two runs)

> Current `run` uses a synthetic metrics generator as placeholder.
> Next step is wiring real S3 operations (PUT/GET/DELETE) via AWS SDK v2.

## Build binary

```bash
mkdir -p bin
go build -o bin/s3bench ./cmd/s3bench
```

Run from source (quick dev):

```bash
go run ./cmd/s3bench validate -c bench.yaml
go run ./cmd/s3bench run -c bench.yaml -o run1.json
go run ./cmd/s3bench run -c bench.yaml -o run2.json
go run ./cmd/s3bench compare run1.json run2.json
```

Run with built binary:

```bash
./bin/s3bench validate -c bench.yaml
./bin/s3bench run -c bench.yaml -o run1.json
./bin/s3bench run -c bench.yaml -o run2.json
./bin/s3bench compare run1.json run2.json
```

## Config reference (bench.yaml)

`bench.yaml` is now documented inline with comments for each field.
Key points:

- `execution.warmup` is excluded from final metrics
- `execution.duration` is the measured phase
- `execution.workers` controls concurrency
- `workload.mode=mixed` requires `mix.put + mix.get + mix.delete = 100`

## Expected output

### `validate`

```text
ok
```

### `run`

```text
done: run1.json (+ run1.csv)
```

Creates:
- `run1.json` (structured result)
- `run1.csv` (tabular metrics)

Example `run1.json` snippet:

```json
{
  "version": "v1-dev",
  "timestamp": "2026-03-12T11:00:00+01:00",
  "duration": "30s",
  "endpoint": "https://s3.example.com",
  "bucket": "perf-test-bucket",
  "metrics": {
    "put": { "ops": 12000, "success": 11820, "errors": 180, "ops_per_sec": 400.0, "mb_per_sec": 7.3 },
    "get": { "ops": 10000, "success": 9850, "errors": 150, "ops_per_sec": 333.33, "mb_per_sec": 6.0 },
    "delete": { "ops": 8000, "success": 7880, "errors": 120, "ops_per_sec": 266.67, "mb_per_sec": 4.5 }
  }
}
```

### `compare`

```text
Compare (B vs A)
- put: ops/s 400.00 -> 400.00 (+0.00)
- delete: ops/s 266.67 -> 266.67 (+0.00)
- get: ops/s 333.33 -> 333.33 (+0.00)
```
