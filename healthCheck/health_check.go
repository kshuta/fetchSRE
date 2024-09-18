package healthCheck

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Endpoint struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

func Run(fname string, logger *slog.Logger) error {

	bdata, err := os.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	logger.Debug("config file:", "data", string(bdata))

	var endpoints []Endpoint
	err = yaml.Unmarshal(bdata, &endpoints)
	if err != nil {
		return fmt.Errorf("unmarshalling config file: %w", err)
	}

	logger.Debug("marshalled data:", "endpoints", endpoints)

	return nil
}
