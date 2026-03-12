package runner

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/indor79/s3bench/internal/config"
	"github.com/indor79/s3bench/internal/s3client"
)

type opAgg struct {
	ops     atomic.Int64
	success atomic.Int64
	errors  atomic.Int64
	bytes   atomic.Int64
	latMu   sync.Mutex
	latency []float64
}

func (a *opAgg) addLatency(ms float64) {
	a.latMu.Lock()
	a.latency = append(a.latency, ms)
	a.latMu.Unlock()
}

type keyPool struct {
	mu   sync.Mutex
	keys []string
}

func (p *keyPool) add(k string) {
	p.mu.Lock()
	p.keys = append(p.keys, k)
	p.mu.Unlock()
}

func (p *keyPool) any() (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.keys) == 0 {
		return "", false
	}
	return p.keys[rand.Intn(len(p.keys))], true
}

func (p *keyPool) popAny() (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.keys) == 0 {
		return "", false
	}
	i := rand.Intn(len(p.keys))
	k := p.keys[i]
	p.keys[i] = p.keys[len(p.keys)-1]
	p.keys = p.keys[:len(p.keys)-1]
	return k, true
}

func Run(c config.Config) (Result, error) {
	rand.Seed(time.Now().UnixNano())

	warmup, _ := time.ParseDuration(c.Execution.Warmup)
	duration, _ := time.ParseDuration(c.Execution.Duration)
	perReqTimeout, _ := time.ParseDuration(c.Execution.PerRequestTimeout)

	ctx := context.Background()
	cli, err := s3client.New(ctx, c)
	if err != nil {
		return Result{}, err
	}

	size := int64(1024 * 1024) // v1 default: 1MiB
	payload := bytes.Repeat([]byte("a"), int(size))

	agg := map[string]*opAgg{
		"put":    {},
		"get":    {},
		"delete": {},
	}

	pool := &keyPool{}
	prefix := strings.Trim(c.Prefix, "/")
	if prefix != "" {
		prefix += "/"
	}
	runID := time.Now().UnixNano()

	chooseOp := func() string {
		switch strings.ToLower(c.Workload.Mode) {
		case "put", "get", "delete":
			return strings.ToLower(c.Workload.Mode)
		case "mixed":
			n := rand.Intn(100) + 1
			if n <= c.Workload.Mix.Put {
				return "put"
			}
			if n <= c.Workload.Mix.Put+c.Workload.Mix.Get {
				return "get"
			}
			return "delete"
		default:
			return "put"
		}
	}

	var runOp func(op string, collect bool)
	runOp = func(op string, collect bool) {
		opCtx, cancel := context.WithTimeout(ctx, perReqTimeout)
		defer cancel()
		start := time.Now()
		var opErr error

		switch op {
		case "put":
			key := fmt.Sprintf("%s%d-%d", prefix, runID, rand.Int63())
			_, opErr = cli.PutObject(opCtx, &s3.PutObjectInput{
				Bucket: &c.Bucket,
				Key:    &key,
				Body:   bytes.NewReader(payload),
			})
			if opErr == nil {
				pool.add(key)
			}
		case "get":
			key, ok := pool.any()
			if !ok {
				runOp("put", collect)
				return
			}
			obj, err := cli.GetObject(opCtx, &s3.GetObjectInput{Bucket: &c.Bucket, Key: &key})
			opErr = err
			if err == nil && obj != nil {
				_ = obj.Body.Close()
			}
		case "delete":
			key, ok := pool.popAny()
			if !ok {
				runOp("put", collect)
				return
			}
			_, opErr = cli.DeleteObject(opCtx, &s3.DeleteObjectInput{Bucket: &c.Bucket, Key: &key})
		}

		if !collect {
			return
		}

		a := agg[op]
		a.ops.Add(1)
		if opErr != nil {
			a.errors.Add(1)
		} else {
			a.success.Add(1)
			a.bytes.Add(size)
		}
		a.addLatency(float64(time.Since(start).Milliseconds()))
	}

	worker := func(stop <-chan struct{}, collect bool, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				runOp(chooseOp(), collect)
			}
		}
	}

	// Warmup phase
	if warmup > 0 {
		stopWarm := make(chan struct{})
		var wg sync.WaitGroup
		for i := 0; i < c.Execution.Workers; i++ {
			wg.Add(1)
			go worker(stopWarm, false, &wg)
		}
		time.Sleep(warmup)
		close(stopWarm)
		wg.Wait()
	}

	// Measure phase
	stopRun := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < c.Execution.Workers; i++ {
		wg.Add(1)
		go worker(stopRun, true, &wg)
	}
	time.Sleep(duration)
	close(stopRun)
	wg.Wait()

	toMetrics := func(a *opAgg) OpMetrics {
		ops := a.ops.Load()
		success := a.success.Load()
		errors := a.errors.Load()
		sec := duration.Seconds()
		if sec <= 0 {
			sec = 1
		}

		a.latMu.Lock()
		lat := append([]float64(nil), a.latency...)
		a.latMu.Unlock()
		sort.Float64s(lat)
		pct := func(p float64) float64 {
			if len(lat) == 0 {
				return 0
			}
			i := int((p / 100) * float64(len(lat)-1))
			if i < 0 {
				i = 0
			}
			if i >= len(lat) {
				i = len(lat) - 1
			}
			return lat[i]
		}

		errorRate := 0.0
		if ops > 0 {
			errorRate = (float64(errors) / float64(ops)) * 100
		}

		return OpMetrics{
			Ops:          ops,
			Success:      success,
			Errors:       errors,
			OpsPerSec:    float64(ops) / sec,
			MBPerSec:     (float64(a.bytes.Load()) / (1024 * 1024)) / sec,
			P50Ms:        pct(50),
			P95Ms:        pct(95),
			P99Ms:        pct(99),
			ErrorRatePct: errorRate,
		}
	}

	res := Result{
		Version:   "v1-sdk",
		Timestamp: time.Now(),
		Duration:  c.Execution.Duration,
		Endpoint:  c.Endpoint,
		Bucket:    c.Bucket,
		Metrics: map[string]OpMetrics{
			"put":    toMetrics(agg["put"]),
			"get":    toMetrics(agg["get"]),
			"delete": toMetrics(agg["delete"]),
		},
	}

	// best-effort bucket cleanup marker for eventual consistency awareness
	_, _ = cli.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &c.Bucket,
		Key:    awsStr(fmt.Sprintf("%s%s", prefix, ".s3bench-run-complete")),
		Body:   bytes.NewReader([]byte(time.Now().Format(time.RFC3339))),
		ACL:    types.ObjectCannedACLPrivate,
	})

	return res, nil
}

func awsStr(s string) *string { return &s }
