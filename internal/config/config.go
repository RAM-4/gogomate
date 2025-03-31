package config

import (
	"fmt"
	"os"
)

type Config struct {
	MistralAPIKey      string
	MistralAgentID     string
	MistralAPIEndpoint string
}

func requireEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("required environment variable %q is not set", key)
	}
	return value, nil
}

func Load() (*Config, error) {
	mistralAPIKey, err := requireEnv("MISTRAL_API_KEY")
	if err != nil {
		return nil, err
	}

	mistralAgentID, err := requireEnv("MISTRAL_AGENT_ID")
	if err != nil {
		return nil, err
	}

	return &Config{
		MistralAPIKey:      mistralAPIKey,
		MistralAgentID:     mistralAgentID,
		MistralAPIEndpoint: "https://api.mistral.ai/v1/agents/completions",
	}, nil
}
