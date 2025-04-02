package cli

import (
	"fmt"
	"net/http"
	"net/url"

	"gogomate/internal/client"
	"gogomate/internal/config"

	"github.com/atotto/clipboard"
	"github.com/urfave/cli/v2"
)

type clients struct {
	scraper *client.Scraper
	mistral *client.MistralClient
}

func NewCLI(cfg *config.Config) *cli.App {
	httpClient := &http.Client{}
	c := &clients{
		scraper: client.NewScraper(
			httpClient,
		),
		mistral: client.NewMistralClient(
			httpClient,
			cfg.MistralAPIKey,
			cfg.MistralAgentID,
			cfg.MistralAPIEndpoint,
		),
	}

	return &cli.App{
		Name:        "gogomate",
		Usage:       "AI-powered cover letter generator",
		Description: "Generate a personalized cover letter from a job posting URL",
		Commands: []*cli.Command{
			{
				Name:      "generate",
				Aliases:   []string{"gen"},
				Usage:     "Generate a cover letter from a job posting URL",
				ArgsUsage: "URL",
				Action:    c.generateCoverLetterFromURL,
			},
		},
	}
}

func (c *clients) generateCoverLetterFromURL(ctx *cli.Context) error {
	url, err := c.validateURL(ctx)
	if err != nil {
		return err
	}

	return c.generateCoverLetter(url)
}

func (c *clients) validateURL(ctx *cli.Context) (string, error) {
	if ctx.NArg() < 1 {
		return "", fmt.Errorf("missing URL argument")
	}

	urlStr := ctx.Args().First()
	if _, err := url.ParseRequestURI(urlStr); err != nil {
		return "", fmt.Errorf("invalid URL: %w\nTip: If your URL contains special characters, wrap it in quotes", err)
	}
	return urlStr, nil
}

func (c *clients) generateCoverLetter(urlStr string) error {
	content, err := c.scraper.Content(urlStr)
	if err != nil {
		return fmt.Errorf("error scraping content: %w", err)
	}

	coverLetter, err := c.mistral.GenerateLetter(content)
	if err != nil {
		return fmt.Errorf("error generating cover letter: %w", err)
	}

	fmt.Println(coverLetter)
	return clipboard.WriteAll(coverLetter)
}
