package compare

import (
	"encoding/json"
	"fmt"
	"os"

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
	for op, am := range a.Metrics {
		bm := b.Metrics[op]
		delta := bm.OpsPerSec - am.OpsPerSec
		out += fmt.Sprintf("- %s: ops/s %.2f -> %.2f (%+.2f)\n", op, am.OpsPerSec, bm.OpsPerSec, delta)
	}
	return out
}
