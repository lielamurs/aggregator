package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Body)
}

type HTTPClient struct {
	client *http.Client
	logger *logrus.Logger
}

func NewHTTPClient(timeout time.Duration, logger *logrus.Logger) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

func (c *HTTPClient) PostJSON(ctx context.Context, url string, payload any, response any) error {
	return c.makeJSONRequest(ctx, "POST", url, payload, response)
}

func (c *HTTPClient) GetJSON(ctx context.Context, url string, response any) error {
	return c.makeJSONRequest(ctx, "GET", url, nil, response)
}

func (c *HTTPClient) makeJSONRequest(ctx context.Context, method, url string, payload any, response any) error {
	logger := c.logger.WithFields(logrus.Fields{
		"method": method,
		"url":    url,
	})

	var requestBody io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			logger.WithError(err).Error("Failed to marshal request payload")
			return fmt.Errorf("failed to marshal request payload: %w", err)
		}
		requestBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, requestBody)
	if err != nil {
		logger.WithError(err).Error("Failed to create HTTP request")
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		logger.WithError(err).Error("HTTP request failed")
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
