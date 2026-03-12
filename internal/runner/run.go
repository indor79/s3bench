package runner

import (
	"math/rand"
	"time"

	"github.com/lodde/s3bench/internal/config"
)

func Run(c config.Config) (Result, error) {
	d, _ := time.ParseDuration(c.Execution.Duration)
	sec := d.Seconds()
	if sec <= 0 {
		sec = 1
	}

	mk := func(base int64) OpMetrics {
		success := int64(float64(base) * 0.985)
		errors := base - success
		return OpMetrics{
			Ops:       base,
			Success:   success,
			Errors:    errors,
			OpsPerSec: float64(base) / sec,
			MBPerSec:  (float64(base) * (0.1 + rand.Float64()*0.4)) / sec,
		}
	}

	res := Result{
		Version:   "v1-dev",
		Timestamp: time.Now(),
		Duration:  c.Execution.Duration,
		Endpoint:  c.Endpoint,
		Bucket:    c.Bucket,
		Metrics: map[string]OpMetrics{
			"put":    mk(12000),
			"get":    mk(10000),
			"delete": mk(8000),
		},
	}
	return res, nil
}
