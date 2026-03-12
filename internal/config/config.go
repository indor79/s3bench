package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Endpoint string `yaml:"endpoint"`
	Region   string `yaml:"region"`
	Bucket   string `yaml:"bucket"`
	Prefix   string `yaml:"prefix"`

	Auth struct {
		AccessKeyEnv   string `yaml:"access_key_env"`
		SecretKeyEnv   string `yaml:"secret_key_env"`
		SessionTokenEnv string `yaml:"session_token_env"`
	} `yaml:"auth"`

	Execution struct {
		Warmup            string `yaml:"warmup"`
		Duration          string `yaml:"duration"`
		Workers           int    `yaml:"workers"`
		PerRequestTimeout string `yaml:"per_request_timeout"`
	} `yaml:"execution"`

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
	return c, Validate(c)
}

func Validate(c Config) error {
	if strings.TrimSpace(c.Endpoint) == "" {
		return errors.New("endpoint is required")
	}
	if strings.TrimSpace(c.Bucket) == "" {
		return errors.New("bucket is required")
	}
	if strings.TrimSpace(c.Workload.Mode) == "" {
		return errors.New("workload.mode is required")
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
