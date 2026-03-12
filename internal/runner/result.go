package runner

import "time"

type OpMetrics struct {
	Ops       int64   `json:"ops"`
	Success   int64   `json:"success"`
	Errors    int64   `json:"errors"`
	OpsPerSec float64 `json:"ops_per_sec"`
	MBPerSec  float64 `json:"mb_per_sec"`
}

type Result struct {
	Version   string               `json:"version"`
	Timestamp time.Time            `json:"timestamp"`
	Duration  string               `json:"duration"`
	Endpoint  string               `json:"endpoint"`
	Bucket    string               `json:"bucket"`
	Metrics   map[string]OpMetrics `json:"metrics"`
}
