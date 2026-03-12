package compare

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/indor79/s3bench/internal/runner"
)

func Load(path string) (runner.Result, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return runner.Result{}, err
	}
	var r runner.Result
	if err := json.Unmarshal(b, &r); err != nil {
		return runner.Result{}, err
	}
	return r, nil
}

func Summary(a, b runner.Result) string {
	out := "Compare (B vs A)\n"
	ops := make([]string, 0, len(a.Metrics))
	for op := range a.Metrics {
		ops = append(ops, op)
	}
	sort.Strings(ops)

	pct := func(old, new float64) float64 {
		if old == 0 {
			return 0
		}
		return ((new - old) / old) * 100
	}

	for _, op := range ops {
		am := a.Metrics[op]
		bm := b.Metrics[op]
		out += fmt.Sprintf("- %s\n", op)
		out += fmt.Sprintf("  ops/s: %.2f -> %.2f (%+.2f%%)\n", am.OpsPerSec, bm.OpsPerSec, pct(am.OpsPerSec, bm.OpsPerSec))
		out += fmt.Sprintf("  mb/s:  %.2f -> %.2f (%+.2f%%)\n", am.MBPerSec, bm.MBPerSec, pct(am.MBPerSec, bm.MBPerSec))
		out += fmt.Sprintf("  p95ms: %.2f -> %.2f (%+.2f%%)\n", am.P95Ms, bm.P95Ms, pct(am.P95Ms, bm.P95Ms))
		out += fmt.Sprintf("  err%%:  %.2f -> %.2f (%+.2f%%)\n", am.ErrorRatePct, bm.ErrorRatePct, pct(am.ErrorRatePct, bm.ErrorRatePct))
	}
	return out
}
