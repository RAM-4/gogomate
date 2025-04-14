package config

import (
	"fmt"
	"os"
)

const (
	endpoint   = "https://api.mistral.ai/v1/agents/completions"
	apiKeyEnv  = "MISTRAL_API_KEY"
	agentIDEnv = "MISTRAL_AGENT_ID"
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
	apiKey, err := requireEnv(apiKeyEnv)
	if err != nil {
		return nil, err
	}

	agentID, err := requireEnv(agentIDEnv)
	if err != nil {
		return nil, err
	}

	return &Config{
		MistralAPIKey:      apiKey,
		MistralAgentID:     agentID,
		MistralAPIEndpoint: endpoint,
	}, nil
}
