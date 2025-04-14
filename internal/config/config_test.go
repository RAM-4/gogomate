package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name: "all keys set",
			env: map[string]string{
				apiKeyEnv:  "test-key",
				agentIDEnv: "test-agent",
			},
			want: &Config{
				MistralAPIKey:      "test-key",
				MistralAgentID:     "test-agent",
				MistralAPIEndpoint: endpoint,
			},
		},
		{
			name: "missing api key",
			env: map[string]string{
				agentIDEnv: "test-agent",
			},
			wantErr: true,
		},
		{
			name: "missing agent id",
			env: map[string]string{
				apiKeyEnv: "test-key",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			t.Cleanup(func() {
				os.Unsetenv(apiKeyEnv)
				os.Unsetenv(agentIDEnv)
			})

			got, err := Load()
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
