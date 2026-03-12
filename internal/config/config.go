package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/indor79/s3bench/internal/util"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Endpoint string `yaml:"endpoint"`
	Region   string `yaml:"region"`
	Bucket   string `yaml:"bucket"`
	Prefix   string `yaml:"prefix"`

	Auth struct {
		AccessKeyEnv    string `yaml:"access_key_env"`
		SecretKeyEnv    string `yaml:"secret_key_env"`
		SessionTokenEnv string `yaml:"session_token_env"`
	} `yaml:"auth"`

	Execution struct {
		Warmup            string `yaml:"warmup"`
		Duration          string `yaml:"duration"`
		Workers           int    `yaml:"workers"`
		PerRequestTimeout string `yaml:"per_request_timeout"`
	} `yaml:"execution"`

	Dataset struct {
		ObjectSizes    []string `yaml:"object_sizes"`
		PrefillObjects int      `yaml:"prefill_objects"`
		KeyMode        string   `yaml:"key_mode"`
	} `yaml:"dataset"`

	Workload struct {
		Mode string `yaml:"mode"`
		Mix  struct {
			Put    int `yaml:"put"`
			Get    int `yaml:"get"`
			Delete int `yaml:"delete"`
		} `yaml:"mix"`
	} `yaml:"workload"`
}

func Load(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return Config{}, err
	}
	if err := Validate(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}

func Validate(c *Config) error {
	if strings.TrimSpace(c.Endpoint) == "" {
		return errors.New("endpoint is required")
	}
	if strings.TrimSpace(c.Bucket) == "" {
		return errors.New("bucket is required")
	}
	if len(c.Dataset.ObjectSizes) == 0 {
		c.Dataset.ObjectSizes = []string{"1MiB"}
	}
	if c.Dataset.PrefillObjects < 0 {
		return errors.New("dataset.prefill_objects must be >= 0")
	}
	if strings.TrimSpace(c.Dataset.KeyMode) == "" {
		c.Dataset.KeyMode = "deterministic"
	}
	km := strings.ToLower(strings.TrimSpace(c.Dataset.KeyMode))
	if km != "deterministic" && km != "random" {
		return errors.New("dataset.key_mode must be deterministic|random")
	}
	for _, s := range c.Dataset.ObjectSizes {
		if _, err := util.ParseSize(s); err != nil {
			return fmt.Errorf("dataset.object_sizes invalid: %w", err)
		}
	}

	if strings.TrimSpace(c.Workload.Mode) == "" {
		return errors.New("workload.mode is required")
	}
	mode := strings.ToLower(strings.TrimSpace(c.Workload.Mode))
	if mode != "put" && mode != "get" && mode != "delete" && mode != "mixed" {
		return errors.New("workload.mode must be one of put|get|delete|mixed")
	}
	if c.Execution.Workers <= 0 {
		return errors.New("execution.workers must be > 0")
	}
	if _, err := time.ParseDuration(c.Execution.Warmup); err != nil {
		return fmt.Errorf("execution.warmup invalid: %w", err)
	}
	if _, err := time.ParseDuration(c.Execution.Duration); err != nil {
		return fmt.Errorf("execution.duration invalid: %w", err)
	}
	if _, err := time.ParseDuration(c.Execution.PerRequestTimeout); err != nil {
		return fmt.Errorf("execution.per_request_timeout invalid: %w", err)
	}
	if c.Workload.Mode == "mixed" {
		total := c.Workload.Mix.Put + c.Workload.Mix.Get + c.Workload.Mix.Delete
		if total != 100 {
			return errors.New("workload.mix must sum to 100")
		}
	}
	return nil
}
