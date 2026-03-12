package runner

import (
	"fmt"

	"github.com/indor79/s3bench/internal/config"
)

func CheckThresholds(cfg config.Config, r Result) []string {
	var fails []string
	for op, m := range r.Metrics {
		if cfg.Thresholds.P95MsMax > 0 && m.P95Ms > cfg.Thresholds.P95MsMax {
			fails = append(fails, fmt.Sprintf("%s p95_ms %.2f > %.2f", op, m.P95Ms, cfg.Thresholds.P95MsMax))
		}
		if cfg.Thresholds.ErrorRateMaxPct > 0 && m.ErrorRatePct > cfg.Thresholds.ErrorRateMaxPct {
			fails = append(fails, fmt.Sprintf("%s error_rate_pct %.2f > %.2f", op, m.ErrorRatePct, cfg.Thresholds.ErrorRateMaxPct))
		}
		if cfg.Thresholds.MinThroughputMBPS > 0 && m.MBPerSec < cfg.Thresholds.MinThroughputMBPS {
			fails = append(fails, fmt.Sprintf("%s mb_per_sec %.2f < %.2f", op, m.MBPerSec, cfg.Thresholds.MinThroughputMBPS))
		}
	}
	return fails
}
