package healthCheck

import (
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const INTERVAL = 3

type Endpoint struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

type Counter struct {
	success, fail int
}

func Run(fname string, sigChan chan os.Signal, logger *slog.Logger) error {
	logger.Info("running health check", "filename", fname)

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

	client := http.Client{
		Timeout: 500 * time.Millisecond,
	}

	domain2Counter, host2Domain, err := getDomainMap(endpoints)
	if err != nil {
		return fmt.Errorf("getting domain map: err=%w", err)
	}

	logger.Info("got domain map", "host2Domain", host2Domain, "domain2counter", domain2Counter)

	timer := time.NewTimer(0)

	for {
		select {
		case <-timer.C:
			logger.Info("checking health")
			err := checkHealth(client, endpoints, domain2Counter, host2Domain, logger)
			if err != nil {
				return fmt.Errorf("checking health: %w", err)
			}
			timer.Reset(INTERVAL * time.Second)
		case <-sigChan:
			logger.Info("received signal, exiting")
			return nil
		}
	}
}

func checkHealth(
	client http.Client,
	endpoints []Endpoint,
	domain2Counter map[string]*Counter,
	host2Domain map[string]string,
	logger *slog.Logger,
) error {
	for _, ep := range endpoints {
		req, err := http.NewRequest(ep.Method, ep.URL, strings.NewReader(ep.Body))
		if err != nil {
			return fmt.Errorf("creating request: endpoint=%+v, err=%w", ep, err)
		}

		logger.Debug("generated request:", "req", req)

		for key, value := range ep.Headers {
			req.Header.Set(key, value)
		}
		res, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("sending request: req=%+v, err=%w", req, err)
		}

		logger.Debug("response from endpoint:", "res", res)

		if checkStatusCode(res.StatusCode) {
			domain2Counter[host2Domain[ep.URL]].success++
		} else {
			domain2Counter[host2Domain[ep.URL]].fail++
		}
	}

	for domain, counter := range domain2Counter {
		up, down := float64(counter.success), float64(counter.fail)
		fmt.Printf("%s has %.0f%% availability percentage.\n", domain, math.Round(100*(up/(up+down))))
	}

	return nil
}

func checkStatusCode(statusCode int) bool {
	return 200 <= statusCode && statusCode < 300
}

func getDomainMap(endpoints []Endpoint) (map[string]*Counter, map[string]string, error) {
	storage := make(map[string]*Counter)
	hostMap := make(map[string]string)
	for _, ep := range endpoints {
		parsedURL, err := url.Parse(ep.URL)
		if err != nil {
			return nil, nil, fmt.Errorf("parsing URL: url=%s, err=%w", ep.URL, err)
		}
		if _, ok := storage[parsedURL.Host]; !ok {
			storage[parsedURL.Host] = &Counter{}
		}
		hostMap[ep.URL] = parsedURL.Host
	}

	return storage, hostMap, nil
}
