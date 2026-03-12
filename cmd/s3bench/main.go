package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/indor79/s3bench/internal/compare"
	"github.com/indor79/s3bench/internal/config"
	"github.com/indor79/s3bench/internal/report"
	"github.com/indor79/s3bench/internal/runner"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "validate":
		cmdValidate(os.Args[2:])
	case "run":
		cmdRun(os.Args[2:])
	case "compare":
		cmdCompare(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func cmdValidate(args []string) {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	cfgPath := fs.String("c", "bench.yaml", "config path")
	_ = fs.Parse(args)
	_, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Println("invalid:", err)
		os.Exit(2)
	}
	fmt.Println("ok")
}

func cmdRun(args []string) {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	cfgPath := fs.String("c", "bench.yaml", "config path")
	out := fs.String("o", "result.json", "output json")
	_ = fs.Parse(args)

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Println("config error:", err)
		os.Exit(2)
	}
	res, err := runner.Run(cfg)
	if err != nil {
		fmt.Println("run error:", err)
		os.Exit(4)
	}

	report.PrintSummary(res)

	if err := report.WriteJSON(*out, res); err != nil {
		fmt.Println("write json error:", err)
		os.Exit(4)
	}
	csvPath := strings.TrimSuffix(*out, filepath.Ext(*out)) + ".csv"
	if err := report.WriteCSV(csvPath, res); err != nil {
		fmt.Println("write csv error:", err)
		os.Exit(4)
	}
	fmt.Printf("done: %s (+ %s)\n", *out, csvPath)

	fails := runner.CheckThresholds(cfg, res)
	if len(fails) > 0 {
		fmt.Println("thresholds: FAIL")
		for _, f := range fails {
			fmt.Println("-", f)
		}
		os.Exit(3)
	}
	fmt.Println("thresholds: PASS")
}

func cmdCompare(args []string) {
	if len(args) != 2 {
		fmt.Println("usage: s3bench compare run-a.json run-b.json")
		os.Exit(2)
	}
	a, err := compare.Load(args[0])
	if err != nil {
		fmt.Println("compare error:", err)
		os.Exit(2)
	}
	b, err := compare.Load(args[1])
	if err != nil {
		fmt.Println("compare error:", err)
		os.Exit(2)
	}
	fmt.Print(compare.Summary(a, b))
}

func usage() {
	fmt.Println("s3bench <validate|run|compare>")
}
