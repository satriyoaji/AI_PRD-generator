package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type AIService struct {
	anthropicAPIKey string
	openaiAPIKey   string
	httpClient     *http.Client
}

func NewAIService(anthropicKey, openaiKey string) *AIService {
	return&AIService{
		anthropicAPIKey: anthropicKey,
		openaiAPIKey:   openaiKey,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicRequest struct {
	Model       string             `json:"model"`
	MaxTokens   int                `json:"max_tokens"`
	System      string             `json:"system"`
	Messages    []AnthropicMessage `json:"messages"`
}

type AnthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

type OpenAIRequest struct {
	Model    string                      `json:"model"`
	Messages []map[string]interface{}     `json:"messages"`
	MaxTokens int                        `json:"max_tokens"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (s *AIService) GeneratePRD(transcript string) (string, error) {
	prompt := fmt.Sprintf(`You are a product manager assistant. Based on the following meeting transcript, generate a comprehensive Product Requirements Document (PRD).

Meeting Transcript:
%s

Generate a PRD with the following sections:
1. Executive Summary
2. Goals and Objectives
3. User Personas (at least 2-3 personas)
4. Feature Requirements (as a numbered list)
5. Non-Functional Requirements
6. Success Metrics
7. Timeline/Roadmap (if mentioned)

Format the output in clean Markdown. Do not include any preamble or explanation, just the PRD content.`, transcript)

	return s.callClaude(prompt)
}

func (s *AIService) GenerateUserStories(transcript string) ([]map[string]interface{}, error) {
	prompt := fmt.Sprintf(`You are a product manager assistant. Based on the following meeting transcript, generate user stories.

Meeting Transcript:
%s

Generate 5-8 user stories in the following JSON format (return ONLY valid JSON, no markdown code blocks):
[
  {
    "id": "US-001",
    "title": "As a [persona], I want [goal] so that [benefit]",
    "description": "Detailed description of the user story",
    "acceptance_criteria": ["Criterion 1", "Criterion 2", "Criterion 3"]
  }
]

Make sure the stories are diverse, covering different user types and workflows mentioned in the transcript. Return ONLY valid JSON array.`, transcript)

	result, err := s.callClaude(prompt)
	if err != nil {
		return nil, err
	}

	// Clean the result - remove markdown code blocks if present
	result = strings.TrimPrefix(result, "```json")
	result = strings.TrimPrefix(result, "```")
	result = strings.TrimSuffix(result, "```")
	result = strings.TrimSpace(result)

	var stories []map[string]interface{}
	if err := json.Unmarshal([]byte(result), &stories); err != nil {
		log.Printf("Failed to parse user stories: %v, raw result: %s", err, result)
		return nil, fmt.Errorf("failed to parse user stories: %w", err)
	}

	return stories, nil
}

func (s *AIService) GenerateRequirements(transcript string) (string, error) {
	prompt := fmt.Sprintf(`You are a product manager assistant. Based on the following meeting transcript, generate functional requirements.

Meeting Transcript:
%s

Generate functional requirements in the following format:
## Functional Requirements

### FR-001: [Requirement Title]
**Description:** Detailed description of the requirement
**Priority:** High/Medium/Low
**Acceptance Criteria:**
- Criterion 1
- Criterion 2
- Criterion 3

### FR-002: [Requirement Title]
... and so on

Generate 6-10 functional requirements covering:
- Core features
- User interactions
- Data handling
- Integration requirements
- Edge cases

Format the output in clean Markdown. Do not include any preamble or explanation.`, transcript)

	return s.callClaude(prompt)
}

func (s *AIService) GenerateWireframe(transcript string) (map[string]string, error) {
	prompt := fmt.Sprintf(`You are a UI/UX designer. Based on the following meeting transcript, generate a UI wireframe description.

Meeting Transcript:
%s

Generate a Mermaid flowchart that shows the key screens and user flows for this product. Focus on:
1. Main navigation structure
2. Key user flows
3. Screen hierarchy

Return your response as a JSON object with this exact format (return ONLY valid JSON, no markdown code blocks):
{
  "mermaid_code": "graph TD\\n  A[Home] --> B[Login]\\n  B --> C[Dashboard]\\n  ...",
  "description": "Description of the wireframe and key screens"
}

Make the Mermaid code valid and renderable. Return ONLY valid JSON.`, transcript)

	result, err := s.callOpenAI(prompt)
	if err != nil {
		return nil, err
	}

	// Clean the result
	result = strings.TrimPrefix(result, "```json")
	result = strings.TrimPrefix(result, "```")
	result = strings.TrimSuffix(result, "```")
	result = strings.TrimSpace(result)

	var wireframe map[string]string
	if err := json.Unmarshal([]byte(result), &wireframe); err != nil {
		log.Printf("Failed to parse wireframe: %v, raw result: %s", err, result)
		return nil, fmt.Errorf("failed to parse wireframe: %w", err)
	}

	return wireframe, nil
}

func (s *AIService) callClaude(prompt string) (string, error) {
	if s.anthropicAPIKey == "" {
		return getMockResponse("claude", prompt), nil
	}

	reqBody := AnthropicRequest{
		Model:     "claude-3-5-sonnet-20241022",
		MaxTokens: 4096,
		System:    "You are a helpful AI assistant that generates structured content.",
		Messages: []AnthropicMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST",
		"https://api.anthropic.com/v1/messages", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.anthropicAPIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Anthropic API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Anthropic API returned status %d: %s", resp.StatusCode, string(body))
	}

	var anthropicResp AnthropicResponse
	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("empty response from Anthropic")
	}

	return anthropicResp.Content[0].Text, nil
}

func (s *AIService) callOpenAI(prompt string) (string, error) {
	if s.openaiAPIKey == "" {
		return getMockResponse("gpt", prompt), nil
	}

	reqBody := OpenAIRequest{
		Model: "gpt-4o",
		Messages: []map[string]interface{}{
			{"role": "system", "content": "You are a helpful AI assistant that generates structured JSON content."},
			{"role": "user", "content": prompt},
		},
		MaxTokens: 4096,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST",
		"https://api.openai.com/v1/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.openaiAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API returned status %d: %s", resp.StatusCode, string(body))
	}

	var openaiResp OpenAIResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from OpenAI")
	}

	return openaiResp.Choices[0].Message.Content, nil
}

func getMockResponse(provider, prompt string) string {
	switch provider {
	case "claude":
		if strings.Contains(strings.ToLower(prompt), "user story") {
			return `[{"id":"US-001","title":"As a user, I want to log in so that I can access my account","description":"User should be able to authenticate using email and password","acceptance_criteria":["User can enter credentials","System validates credentials","User receives feedback on success/failure"]},{"id":"US-002","title":"As a product manager, I want to view my dashboard so that I can see key metrics","description":"Dashboard displays important product metrics and KPIs","acceptance_criteria":["Dashboard loads within 3 seconds","Metrics are displayed in charts","User can filter by date range"]}]`
		}
		return "# Product Requirements Document\n\n## Executive Summary\nThis document outlines the requirements for the product based on the meeting transcript.\n\n## Goals and Objectives\n- Goal 1: Provide users with essential functionality\n- Goal 2: Deliver a seamless user experience\n\n## User Personas\n1. Primary User: Sarah, 35, Product Manager\n2. Secondary User: John, 28, Developer\n\n## Feature Requirements\n1. User authentication\n2. Dashboard view\n3. Data management\n\n## Success Metrics\n- User adoption rate\n- System uptime\n- Customer satisfaction"
	case "gpt":
		return `{"mermaid_code":"graph TD\\n  A[Login] --> B[Dashboard]\\n  B --> C[Settings]\\n  B --> D[Reports]\\n  D --> E[Export Data]","description":"Main user flow showing login, dashboard access, and key features including settings and reports"}`
	}
	return ""
}