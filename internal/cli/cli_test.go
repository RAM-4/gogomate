package cli

import (
	"flag"
	"testing"

	"gogomate/internal/config"

	"github.com/urfave/cli/v2"
)

func TestNewCLI(t *testing.T) {
	cfg := &config.Config{
		APIKey:      "test-key",
		AgentID:     "test-agent",
		APIEndpoint: "https://example.com/api",
	}

	app := NewCLI(cfg)

	if app.Name != "gogomate" {
		t.Errorf("expected app name to be 'gogomate', got %s", app.Name)
	}

	if app.Usage != "AI-powered cover letter generator" {
		t.Errorf("unexpected usage: %s", app.Usage)
	}

	if len(app.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(app.Commands))
	}

	if app.Commands[0].Name != "generate" {
		t.Errorf("expected command name to be 'generate', got %s", app.Commands[0].Name)
	}
}

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantURL   string
		wantComp  string
		wantError bool
	}{
		{
			name:      "valid URL with company",
			args:      []string{"https://example.com/job", "Example Inc"},
			wantURL:   "https://example.com/job",
			wantComp:  "Example Inc",
			wantError: false,
		},
		{
			name:      "valid URL without company",
			args:      []string{"https://example.com/job"},
			wantURL:   "https://example.com/job",
			wantComp:  "",
			wantError: false,
		},
		{
			name:      "no arguments",
			args:      []string{},
			wantError: true,
		},
		{
			name:      "invalid URL",
			args:      []string{"not-a-valid-url"},
			wantError: true,
		},
	}

	c := &clients{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &cli.App{}
			set := flag.NewFlagSet("test", flag.ContinueOnError)
			if len(tt.args) > 0 {
				err := set.Parse(tt.args)
				if err != nil {
					t.Fatalf("failed to parse arguments: %v", err)
				}
			}
			ctx := cli.NewContext(app, set, nil)

			url, company, err := c.validateArgs(ctx)

			if (err != nil) != tt.wantError {
				t.Fatalf("validateArgs() error = %v, wantError = %v", err, tt.wantError)
			}

			if !tt.wantError {
				if url != tt.wantURL {
					t.Errorf("validateArgs() url = %q, want %q", url, tt.wantURL)
				}
				if company != tt.wantComp {
					t.Errorf("validateArgs() company = %q, want %q", company, tt.wantComp)
				}
			}
		})
	}
}
