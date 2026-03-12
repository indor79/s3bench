package report

import (
	"fmt"
	"sort"

	"github.com/indor79/s3bench/internal/runner"
)

func PrintSummary(r runner.Result) {
	fmt.Printf("summary: endpoint=%s bucket=%s duration=%s\n", r.Endpoint, r.Bucket, r.Duration)
	ops := make([]string, 0, len(r.Metrics))
	for op := range r.Metrics {
		ops = append(ops, op)
	}
	sort.Strings(ops)
	for _, op := range ops {
		m := r.Metrics[op]
		fmt.Printf("- %-6s ops/s=%8.2f mb/s=%8.2f p95=%7.2fms err=%6.2f%%\n", op, m.OpsPerSec, m.MBPerSec, m.P95Ms, m.ErrorRatePct)
	}
}
