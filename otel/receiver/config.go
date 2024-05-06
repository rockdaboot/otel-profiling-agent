package otelreceiver

import (
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
)

// Config represents the receiver config settings within the collector's config.yaml
type Config struct {
	Interval  string `mapstructure:"interval"`
	ProjectID string `mapstructure:"projectid"`
}

// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
	interval, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		return fmt.Errorf("interval must be a valid duration string: %s", cfg.Interval)
	}
	if interval.Seconds() < 1 {
		return fmt.Errorf("when defined, the interval has to be set to at least 1 second (1s)")
	}

	if cfg.ProjectID == "" {
		return fmt.Errorf("projectid must be set")
	}
	return nil
}

func createDefaultConfig() component.Config {
	return &Config{
		ProjectID: defaultProjectID,
		Interval:  string(defaultInterval),
	}
}
