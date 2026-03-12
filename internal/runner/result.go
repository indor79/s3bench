package runner

import "time"

type OpMetrics struct {
	Ops          int64   `json:"ops"`
	Success      int64   `json:"success"`
	Errors       int64   `json:"errors"`
	OpsPerSec    float64 `json:"ops_per_sec"`
	MBPerSec     float64 `json:"mb_per_sec"`
	P50Ms        float64 `json:"p50_ms"`
	P95Ms        float64 `json:"p95_ms"`
	P99Ms        float64 `json:"p99_ms"`
	ErrorRatePct float64 `json:"error_rate_pct"`
}

type Result struct {
	Version   string               `json:"version"`
	Timestamp time.Time            `json:"timestamp"`
	Duration  string               `json:"duration"`
	Endpoint  string               `json:"endpoint"`
	Bucket    string               `json:"bucket"`
	Metrics   map[string]OpMetrics `json:"metrics"`
}
