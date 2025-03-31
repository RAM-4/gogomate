package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type requester interface {
	Do(req *http.Request) (*http.Response, error)
}

type MistralClient struct {
	http     requester
	apiKey   string
	endpoint string
	agentID  string
}

type agentResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewMistralClient(httpClient requester, apiKey, agentID, endpoint string) *MistralClient {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &MistralClient{
		http:     httpClient,
		apiKey:   apiKey,
		endpoint: endpoint,
		agentID:  agentID,
	}
}

func (c *MistralClient) GenerateLetter(content string) (string, error) {
	body, err := c.buildPayload(content)
	if err != nil {
		return "", fmt.Errorf("build payload: %w", err)
	}

	resp, err := c.sendRequest(body)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}

	result, err := c.decode(resp)
	if err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}

func (c *MistralClient) buildPayload(content string) ([]byte, error) {
	payload := map[string]any{
		"agent_id": c.agentID,
		"messages": []map[string]string{
			{"role": "user", "content": content},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}
	return data, nil
}

func (c *MistralClient) sendRequest(payload []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return responseBody, nil
}

func (c *MistralClient) decode(body []byte) (string, error) {
	var resp agentResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return resp.Choices[0].Message.Content, nil
}
