package check

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Target is one thing to check.
type Target struct {
	Name    string `yaml:"name"`
	Kind    string `yaml:"kind"` // "tcp", "http" or "dns"
	Address string `yaml:"address"`
}

// Config is the whole targets.yaml file.
type Config struct {
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
	Targets  []Target      `yaml:"targets"`
}

// Load reads and validates a targets file.
func Load(path string) (Config, error) {
	var cfg Config
	// TODO 1: data, err := os.ReadFile(path); wrap errors with fmt.Errorf("read %s: %w", path, err)
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("Error parsing yaml: %w", err)
	}
	// TODO 2: yaml.Unmarshal(data, &cfg)
	// TODO 3: defaults: Interval 30*time.Second, Timeout 5*time.Second when zero
	// TODO 4: error if there are no targets, or a Kind is not tcp/http/dns (a switch works well)
	if cfg.Interval == 0 {
		cfg.Interval = 30 * time.Second
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}
	if len(cfg.Targets) == 0 {
		return cfg, fmt.Errorf("Target List Empty")
	}
	for _, t := range cfg.Targets {
		switch t.Kind {
		case "http", "tcp", "dns":
			continue
		default:
			return cfg, fmt.Errorf("Invalid kind %q", t.Kind)
		}
	}

	return cfg, nil
}
