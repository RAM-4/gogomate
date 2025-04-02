# Gogomate

> A CLI tool that leverages AI to generate personalized cover letters from job posting URLs.

## 🛠️ Features

- Scrapes job posting content from provided URLs
- Generates customized cover letters using Mistral AI
- Automatically saves generated letters to a local directory
- Copies the generated letter to your clipboard
- Simple command-line interface

## 📋 Prerequisites

- Go 1.24.1 or later
- A Mistral AI API key
- A Mistral AI agent ID configured for cover letter generation

## 🚀 Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/gogomate.git
   cd gogomate
   ```

2. Build the binary:
   ```bash
   go build ./cmd/gogomate
   ```

3. (Optional) Move the binary to your PATH:
   ```bash
   mv gogomate /usr/local/bin/
   ```

## ⚙️ Configuration

The tool requires two environment variables:

```bash
export MISTRAL_API_KEY='your-mistral-api-key'
export MISTRAL_AGENT_ID='your-mistral-agent-id'
```

You can add these to your shell's configuration file (e.g., `~/.bashrc` or `~/.zshrc`) for persistence.

## 📖 Usage

### Basic Usage

Generate a cover letter from a job posting URL:

```bash
gogomate generate "https://example.com/job-posting"
```

### With Company Name

Include a company name to better organize saved letters:

```bash
gogomate generate "https://example.com/job-posting" "Example Corp"
```

### Using the Short Command

You can also use the shorter `gen` alias:

```bash
gogomate gen "https://example.com/job-posting" "Example Corp"
```

## 📂 Output

- Generated cover letters are saved in the `letters/` directory
- File naming format: `{company_name}_{timestamp}.txt` or `letter_{timestamp}.txt`
- The generated letter is automatically copied to your clipboard

## 📝 Notes

- If the URL contains special characters, wrap it in quotes
- The tool creates a `letters` directory in your current working directory
- Make sure you have write permissions in the current directory
