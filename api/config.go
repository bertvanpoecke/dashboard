package api

import "github.com/koding/multiconfig"

type config struct {
	Port     int    `default:"8080"`
	LogLevel string `default:"info"`
}

func loadConfiguration() (*config, error) {
	loader := multiconfig.New()
	config := &config{}

	if err := loader.Load(config); err != nil {
		return nil, err
	}

	if err := loader.Validate(config); err != nil {
		return nil, err
	}

	return config, nil
}
