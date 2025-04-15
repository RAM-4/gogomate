package cli

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"gogomate/internal/client"
	"gogomate/internal/config"

	"github.com/atotto/clipboard"
	"github.com/briandowns/spinner"
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
			cfg.APIKey,
			cfg.AgentID,
			cfg.APIEndpoint,
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
				ArgsUsage: "URL [COMPANY_NAME]",
				Action:    c.generateCoverLetterFromURL,
			},
		},
	}
}

func (c *clients) generateCoverLetterFromURL(ctx *cli.Context) error {
	url, company, err := c.validateArgs(ctx)
	if err != nil {
		return err
	}

	return c.generateCoverLetter(url, company)
}

func (c *clients) validateArgs(ctx *cli.Context) (string, string, error) {
	if ctx.NArg() < 1 {
		return "", "", fmt.Errorf("missing URL argument")
	}

	urlStr := ctx.Args().First()
	if _, err := url.ParseRequestURI(urlStr); err != nil {
		return "", "", fmt.Errorf("invalid URL: %w\nTip: If your URL contains special characters, wrap it in quotes", err)
	}

	company := ctx.Args().Get(1)

	return urlStr, company, nil
}

func (c *clients) generateCoverLetter(urlStr, company string) error {
	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	s.Suffix = "  Gogo mate ! Generating cover letter"
	s.Start()

	content, err := c.scraper.Content(urlStr)
	if err != nil {
		return fmt.Errorf("error scraping content: %w", err)
	}

	coverLetter, err := c.mistral.GenerateLetter(content)
	if err != nil {
		return fmt.Errorf("error generating cover letter: %w", err)
	}

	s.Stop()
	if err := saveCoverLetter(coverLetter, company); err != nil {
		return fmt.Errorf("error saving cover letter: %w", err)
	}

	fmt.Println("Copied to clipboard and ready to go mate 🙌")
	return clipboard.WriteAll(coverLetter)
}

func saveCoverLetter(result, company string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	folderPath := filepath.Join(homeDir, ".gogomate", "data")
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	var filename string
	if company != "" {
		filename = fmt.Sprintf("%s_%s.txt", company, timestamp)
	} else {
		filename = fmt.Sprintf("letter_%s.txt", timestamp)
	}

	filePath := filepath.Join(folderPath, filename)

	if err := os.WriteFile(filePath, []byte(result), 0600); err != nil {
		return err
	}

	fmt.Printf("Letter saved to: %s\n", filePath)
	return nil
}
