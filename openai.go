package computeruse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Request represents the structure for sending requests to the OpenAI API
type Request struct {
	Model              string  `json:"model"`
	Input              []Input `json:"input"`
	Text               *Text   `json:"text,omitempty"`
	Tools              []Tool  `json:"tools,omitempty"`
	Temperature        float64 `json:"temperature,omitempty"`
	MaxOutputTokens    int     `json:"max_output_tokens,omitempty"`
	TopP               float64 `json:"top_p,omitempty"`
	Stream             bool    `json:"stream,omitempty"`
	Store              bool    `json:"store,omitempty"`
	Reasoning          any     `json:"reasoning,omitempty"`
	Truncation         string  `json:"truncation,omitempty"`
	PreviousResponseID string  `json:"previous_response_id,omitempty"`
}

// Input represents an input message in the request
type Input struct {
	Type                     string          `json:"type,omitempty"`
	CallID                   string          `json:"call_id,omitempty"`
	Output                   *ComputerOutput `json:"output,omitempty"`
	Role                     string          `json:"role,omitempty"`
	Content                  string          `json:"content,omitempty"`
	AcknowledgedSafetyChecks []SafetyCheck   `json:"acknowledged_safety_checks,omitempty"`
}

// ComputerOutput represents computer output data in the API interaction
type ComputerOutput struct {
	Type       string `json:"type"`
	ImageURL   string `json:"image_url"`
	CurrentURL string `json:"current_url"`
}

// Text represents text format configuration
type Text struct {
	Format Format `json:"format"`
}

// Format represents the format specification for text
type Format struct {
	Type   string         `json:"type"`
	Name   string         `json:"name,omitempty"`
	Strict bool           `json:"strict,omitempty"`
	Schema map[string]any `json:"schema,omitempty"`
}

// Response represents the structure for storing responses from the OpenAI API
type Response struct {
	ID                 string         `json:"id"`
	Object             string         `json:"object"`
	CreatedAt          int            `json:"created_at"`
	Status             string         `json:"status"`
	Error              any            `json:"error"`
	IncompleteDetails  any            `json:"incomplete_details"`
	Instructions       any            `json:"instructions"`
	MaxOutputTokens    any            `json:"max_output_tokens"`
	Model              string         `json:"model"`
	Output             []OutputItem   `json:"output"`
	ParallelToolCalls  bool           `json:"parallel_tool_calls"`
	PreviousResponseID string         `json:"previous_response_id"`
	Reasoning          ReasoningInfo  `json:"reasoning"`
	Store              bool           `json:"store"`
	Temperature        float64        `json:"temperature"`
	Text               TextInfo       `json:"text"`
	ToolChoice         string         `json:"tool_choice"`
	Tools              []any          `json:"tools"`
	TopP               float64        `json:"top_p"`
	Truncation         string         `json:"truncation"`
	Usage              UsageInfo      `json:"usage"`
	User               string         `json:"user"`
	Metadata           map[string]any `json:"metadata"`
}

// OutputItem represents an output item in the API response
type OutputItem struct {
	Type                string        `json:"type"`
	ID                  string        `json:"id,omitempty"`
	CallID              string        `json:"call_id,omitempty"`
	Status              string        `json:"status,omitempty"`
	Action              *Action       `json:"action,omitempty"`
	Role                string        `json:"role,omitempty"`
	Content             []any         `json:"content,omitempty"`
	PendingSafetyChecks []SafetyCheck `json:"pending_safety_checks,omitempty"`
}

// SafetyCheck represents a safety check in the API response
type SafetyCheck struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Action represents an action in the API response
type Action struct {
	Type    string   `json:"type"`
	Keys    []string `json:"keys,omitempty"`
	Button  string   `json:"button,omitempty"`
	Text    string   `json:"text,omitempty"`
	X       int      `json:"x,omitempty"`
	Y       int      `json:"y,omitempty"`
	ScrollX int      `json:"scroll_x,omitempty"`
	ScrollY int      `json:"scroll_y,omitempty"`
}

// Key represents a key-value pair
type Key struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ContentItem represents a content item in the response
type ContentItem struct {
	Type        string        `json:"type"`
	Text        string        `json:"text"`
	Annotations []interface{} `json:"annotations"`
}

// ReasoningInfo represents reasoning information in the response
type ReasoningInfo struct {
	Effort  interface{} `json:"effort"`
	Summary interface{} `json:"summary"`
}

// TextInfo represents text information in the response
type TextInfo struct {
	Format TextFormat `json:"format"`
}

// TextFormat represents the format of text in the response
type TextFormat struct {
	Type   string      `json:"type,omitempty"`
	Name   string      `json:"name,omitempty"`
	Strict bool        `json:"strict,omitempty"`
	Schema interface{} `json:"schema,omitempty"`
}

// UsageInfo represents usage information in the response
type UsageInfo struct {
	InputTokens         int                 `json:"input_tokens"`
	InputTokensDetails  InputTokensDetails  `json:"input_tokens_details"`
	OutputTokens        int                 `json:"output_tokens"`
	OutputTokensDetails OutputTokensDetails `json:"output_tokens_details"`
	TotalTokens         int                 `json:"total_tokens"`
}

// InputTokensDetails represents details about input tokens
type InputTokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

// OutputTokensDetails represents details about output tokens
type OutputTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

// Tool represents a tool configuration for the API
type Tool struct {
	Type          string `json:"type"`
	DisplayWidth  int    `json:"display_width"`
	DisplayHeight int    `json:"display_height"`
	Environment   string `json:"environment"`
}

// Responses sends a request to the OpenAI API and retrieves the response
// Parameters:
// - model: The model name to use (e.g., "gpt-4o")
// - responseID: Previous response ID for conversation continuity
// - input: Array of input messages
func Responses(model string, responseID string, input []Input) (*Response, error) {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	request := Request{
		Model:              model,
		Input:              input,
		PreviousResponseID: responseID,
		Truncation:         "auto",
	}

	request.Tools = []Tool{
		{
			Type:          "computer-preview",
			DisplayWidth:  1024,
			DisplayHeight: 768,
			Environment:   "browser",
		},
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/responses", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Return error if status code is not 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// NewComputerMessage creates a new user message with the given text
func NewComputerMessage(text string) any {
	return struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{
		Role:    "user",
		Content: text,
	}
}
