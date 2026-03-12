package report

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/lodde/s3bench/internal/runner"
)

func WriteJSON(path string, r runner.Result) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func WriteCSV(path string, r runner.Result) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	_ = w.Write([]string{"operation", "ops", "success", "errors", "ops_per_sec", "mb_per_sec", "p50_ms", "p95_ms", "p99_ms", "error_rate_pct"})
	keys := make([]string, 0, len(r.Metrics))
	for k := range r.Metrics {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, op := range keys {
		m := r.Metrics[op]
		_ = w.Write([]string{
			op,
			fmt.Sprintf("%d", m.Ops),
			fmt.Sprintf("%d", m.Success),
			fmt.Sprintf("%d", m.Errors),
			fmt.Sprintf("%.2f", m.OpsPerSec),
			fmt.Sprintf("%.2f", m.MBPerSec),
			fmt.Sprintf("%.2f", m.P50Ms),
			fmt.Sprintf("%.2f", m.P95Ms),
			fmt.Sprintf("%.2f", m.P99Ms),
			fmt.Sprintf("%.2f", m.ErrorRatePct),
		})
	}
	return w.Error()
}
