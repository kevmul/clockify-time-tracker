// Provides a reusable HTTP client for making Clockify API requests
package clockify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Base URL for all Clockify API requests
const baseURL = "https://api.clockify.me/api/v1"

// Client handles all HTTP interactions with the Clockify API
// It stores the API key and reuses an HTTP client for efficiency
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates and returns a new Clockify API client
// This is the constructor function - always use this to create clients
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{}, // Standard HTTP client
	}
}

// doRequest is a private helper method that performs HTTP requests
// It handles:
// - JSON marshaling of request bodies
// - Adding authentication headers
// - Error handling for non-2xx responses
func (c *Client) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader

	// If we have a body, marshal it to JSON
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	// Create the HTTP request
	req, err := http.NewRequest(method, baseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add required headers
	req.Header.Set("X-Api-Key", c.apiKey) // Clockify uses this header for auth
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close() // Always close the response body

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// get performs a GET request - convenience wrapper around doRequest
func (c *Client) Get(endpoint string) ([]byte, error) {
	return c.doRequest("GET", endpoint, nil)
}

// post performs a POST request - convenience wrapper around doRequest
func (c *Client) Post(endpoint string, body interface{}) ([]byte, error) {
	return c.doRequest("POST", endpoint, body)
}
