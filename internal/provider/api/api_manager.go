// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ErrResourceNotFound is returned when the API returns 500 with a message
// indicating the requested entity does not exist (e.g. out-of-band deletion).
// Callers can use errors.Is(err, api.ErrResourceNotFound) to detect and handle
// by removing the resource from state.
var ErrResourceNotFound = errors.New("resource not found")

// ApiManager manages API interactions with the Commander Service Mode.
// Resources and data sources use this struct to make API calls without
// directly accessing the service mode URL and API key.
type ApiManager struct {
	ServiceModeUrl    string
	ServiceModeApiKey string
	HttpClient        *http.Client
	IsMspAccount      bool          // true for MSP accounts, false for Enterprise accounts
	RequestTimeout    time.Duration // used for HTTP client and for polling command result (default 60s)

	// contextMu serializes all context-changing operations (switch-to-msp / switch-to-mc).
	// currentContext tracks backend context: "" = MSP, non-empty = managed company name/id.
	// Used to skip redundant switch-to-msp API calls when already in MSP.
	contextMu      sync.Mutex
	currentContext string
}

// LockContext acquires the context mutex. Call UnlockContext when done (e.g. defer).
// Serializes switch-to-msp / switch-to-mc operations so only one runs at a time.
func (a *ApiManager) LockContext() {
	a.contextMu.Lock()
}

// UnlockContext releases the context mutex.
func (a *ApiManager) UnlockContext() {
	a.contextMu.Unlock()
}

// GetCurrentContext returns the tracked backend context: "" = MSP, non-empty = managed company.
func (a *ApiManager) GetCurrentContext() string {
	return a.currentContext
}

// SetCurrentContext sets the tracked backend context. Call with "" for MSP.
func (a *ApiManager) SetCurrentContext(s string) {
	a.currentContext = s
}

// SubmitRequestResponse represents the response structure when we submit new requests.
type SubmitRequestResponse struct {
	Success   bool   `json:"success"`
	RequestId string `json:"request_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// RequestStatusResponse represents the response structure when we check the status of a request.
type RequestStatusResponse struct {
	Success     bool   `json:"success"`
	RequestId   string `json:"request_id"`
	Command     string `json:"command"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}

// FlexibleMessage is a type that can unmarshal both string and array from JSON
// and always provides a string representation.
type FlexibleMessage string

// UnmarshalJSON implements custom JSON unmarshaling to handle both string and array.
func (fm *FlexibleMessage) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*fm = FlexibleMessage(str)
		return nil
	}

	// If not a string, try to unmarshal as array
	var arr []interface{}
	if err := json.Unmarshal(data, &arr); err == nil {
		// Convert array to string representation
		arrBytes, _ := json.Marshal(arr)
		*fm = FlexibleMessage(string(arrBytes))
		return nil
	}

	// If neither works, try as any and convert to string
	var anyVal interface{}
	if err := json.Unmarshal(data, &anyVal); err == nil {
		anyBytes, _ := json.Marshal(anyVal)
		*fm = FlexibleMessage(string(anyBytes))
		return nil
	}

	// If all else fails, set empty string
	*fm = FlexibleMessage("")
	return nil
}

// String returns the string representation.
func (fm FlexibleMessage) String() string {
	return string(fm)
}

// RequestResultResponse represents the response structure when we get the result of a request.
type RequestResultResponse struct {
	Data    any             `json:"data"`
	Status  string          `json:"status"`
	Message FlexibleMessage `json:"message"`
	Error   string          `json:"error"`
}

// handleAPIErrorResponse checks HTTP status codes and handles error cases.
// For error status codes: reads body and returns error with body content.
// For success status codes (2xx): return nil (caller reads body).
func handleAPIErrorResponse(resp *http.Response) error {
	// Check if it's an error status code first
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read body only for error cases
		bodyBytes, _ := io.ReadAll(resp.Body)

		// Helper function to extract error message from JSON response
		extractErrorMessage := func(body []byte) string {
			var errorResp RequestResultResponse
			if err := json.Unmarshal(body, &errorResp); err == nil {
				if string(errorResp.Message) != "" {
					return string(errorResp.Message)
				}
				if errorResp.Error != "" {
					return errorResp.Error
				}
			}
			return string(body)
		}

		errorMsg := extractErrorMessage(bodyBytes)

		// Handle specific error status codes according to API documentation
		switch resp.StatusCode {
		case http.StatusServiceUnavailable: // 503 - Queue is full
			return fmt.Errorf("queue is full (503): service unavailable, please try again later")

		case http.StatusTooManyRequests: // 429 - Rate limit exceeded
			return fmt.Errorf("rate limit exceeded (429): too many requests, please retry after some time")

		case http.StatusNotFound: // 404 - Request ID not found
			return fmt.Errorf("request id not found (404)")

		case http.StatusInternalServerError: // 500 - Command execution failed
			if strings.Contains(strings.ToLower(errorMsg), "not found") {
				return fmt.Errorf("%w: %s", ErrResourceNotFound, errorMsg)
			}
			return fmt.Errorf("internal server error (500): %s", errorMsg)

		case http.StatusBadRequest: // 400 - Bad request
			return fmt.Errorf("bad request (400): %s", errorMsg)

		default:
			// For any other non-2xx status
			return fmt.Errorf("keeper Security API request failed with status %d: %s", resp.StatusCode, errorMsg)
		}
	}

	// For 2xx status codes, return nil (success - caller reads body)
	return nil
}

// normalizeCommandForShell doubles single quotes in the command so the Commander
// backend can parse it; ” inside single-quoted strings is one literal quote.
// Replaces every ' with ” so apostrophes in names (e.g. O'Brien) are valid.
// escapedApostropheUnit is the POSIX single-quote escaping sequence QuoteShellSingle
// substitutes for each embedded apostrophe (close quote, literal ' via double quotes,
// reopen quote). A span containing this exact sequence was already correctly quoted
// upstream and must be passed through as-is rather than reinterpreted below.
const escapedApostropheUnit = `'"'"'`

func normalizeCommandForShell(command string) string {
	var b strings.Builder
	b.Grow(len(command) + 16)
	i := 0
	for i < len(command) {
		c := command[i]
		if c != '\'' {
			b.WriteByte(c)
			i++
			continue
		}
		// Start of single-quoted span: find content and closing quote
		origStart := i
		start := i + 1
		end := start
		preEscaped := false
		for end < len(command) {
			if command[end] == '\'' {
				// Already-escaped apostrophe unit (produced by QuoteShellSingle): not a
				// closing quote, skip over it and keep scanning for the true end.
				if end+len(escapedApostropheUnit) <= len(command) && command[end:end+len(escapedApostropheUnit)] == escapedApostropheUnit {
					preEscaped = true
					end += len(escapedApostropheUnit)
					continue
				}
				// Closing delimiter if followed by space/punctuation/end; else apostrophe
				if end+1 >= len(command) || command[end+1] == ' ' || command[end+1] == '\'' ||
					(command[end+1] < 'a' || command[end+1] > 'z') && (command[end+1] < 'A' || command[end+1] > 'Z') && (command[end+1] < '0' || command[end+1] > '9') {
					break
				}
			}
			end++
		}
		if preEscaped {
			// The span already contains a valid escaped apostrophe; copying it through
			// unchanged avoids corrupting it by reprocessing it as a naive raw apostrophe.
			b.WriteString(command[origStart : end+1])
			i = end + 1
			continue
		}
		content := command[start:end]
		hasApostrophe := strings.ContainsRune(content, '\'')
		if hasApostrophe {
			// Output as double-quoted; escape any " in content as \"
			b.WriteByte('"')
			for _, r := range content {
				if r == '"' {
					b.WriteString(`\"`)
				} else {
					b.WriteRune(r)
				}
			}
			b.WriteByte('"')
		} else {
			b.WriteByte('\'')
			b.WriteString(content)
			b.WriteByte('\'')
		}
		i = end + 1
	}
	return b.String()
}

// SubmitRequest creates a new API request and returns the parsed response.
// If fileData is non-nil, it is sent in the request body as the "filedata" field.
// fileData is dynamic: pass an object (e.g. map[string]interface{}) for "filedata": {...},
// or a JSON array (e.g. []map[string]interface{}) for "filedata": [...]. It is serialized as-is.
func (a *ApiManager) SubmitRequest(ctx context.Context, command string, fileData interface{}) (*SubmitRequestResponse, error) {
	// Normalize command so single-quoted arguments with apostrophes are valid in shell
	command = normalizeCommandForShell(command)
	// Build request body: command only, or command + filedata
	var jsonBody []byte
	var err error
	if fileData != nil {
		reqBody := map[string]interface{}{
			"command":  command,
			"filedata": fileData,
		}
		jsonBody, err = json.Marshal(reqBody)
	} else {
		reqBody := map[string]string{
			"command": command,
		}
		jsonBody, err = json.Marshal(reqBody)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Build the full URL
	endpoint := a.ServiceModeUrl + "/executecommand-async"

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", a.ServiceModeApiKey)

	// Make the HTTP request using the client
	resp, err := a.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Handle response using common function
	// Expected status: 202 Accepted

	// Check for expected status: 202 Accepted
	if resp.StatusCode == http.StatusAccepted {
		// Read and parse response body (202 Accepted - success)
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var result SubmitRequestResponse
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		return &result, nil
	}

	// Handle all other status codes (errors or unexpected 2xx)
	if apiError := handleAPIErrorResponse(resp); apiError != nil {
		return nil, apiError
	}

	// If we get here, it's a 2xx but not 202 - unexpected
	return nil, fmt.Errorf("unexpected status code: %d, expected 202 Accepted", resp.StatusCode)

}

// RequestStatus checks the status of an API request.
func (a *ApiManager) RequestStatus(ctx context.Context, requestId string) (*RequestStatusResponse, error) {
	// Build the full URL
	endpoint := a.ServiceModeUrl + "/status/" + requestId

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("api-key", a.ServiceModeApiKey)

	// Make the HTTP request using the client
	resp, err := a.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Handle response using common function
	// Expected status: 200 OK
	if apiError := handleAPIErrorResponse(resp); apiError != nil {
		// For 404, provide more context
		if strings.Contains(apiError.Error(), "404") {
			return nil, fmt.Errorf("request ID not found (404): request %s does not exist", requestId)
		}
		return nil, apiError
	}

	// Read and parse response body (200 OK - success)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result RequestStatusResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// RequestResult gets the result of an API request.
func (a *ApiManager) RequestResult(ctx context.Context, requestId string) (*RequestResultResponse, error) {
	// Build the full URL
	endpoint := a.ServiceModeUrl + "/result/" + requestId

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("api-key", a.ServiceModeApiKey)

	// Make the HTTP request using the client
	resp, err := a.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	/*
		THIS /result RETURNS
			- 202 WHEN REQ STATUS IS STILL IN QUEUED/PROCESSING
			- 200 WHEN REQ STATUS IS SUCCESS

	*/

	// Check for 202 (still queued/processing) - this is not an error, just means "not ready yet"
	if resp.StatusCode == http.StatusAccepted {
		return nil, fmt.Errorf("request is still in queued/processing (202)")
	}

	// Handle API errors (404, 500, 503, 429, etc.)
	if apiError := handleAPIErrorResponse(resp); apiError != nil {
		return nil, apiError
	}

	// If we get here, status should be 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, expected 200", resp.StatusCode)
	}

	// Read and parse response body (200 OK - success)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result RequestResultResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ExecuteCommand submits a command and polls for the result.
// This is a convenience function that combines SubmitRequest and PollRequestResult.
// It handles all error checking and returns the final result or an error.
// Optional fileData is sent as the "filedata" field in the request body. It is dynamic:
// pass an object (map) for "filedata": {...} or a JSON array (slice) for "filedata": [...].
func (a *ApiManager) ExecuteCommand(ctx context.Context, command string, errorSummary string, fileData ...interface{}) (*RequestResultResponse, error) {
	var fd interface{}
	if len(fileData) > 0 {
		fd = fileData[0]
	}
	// Submit the request
	submitResp, err := a.SubmitRequest(ctx, command, fd)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorSummary, err)
	}

	// Poll for the result (uses same timeout as HTTP client by default)
	apiResp, err := a.PollRequestResult(ctx, submitResp.RequestId, a.RequestTimeout)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorSummary, err)
	}

	// Check if the result status is success
	if apiResp.Status != "success" {
		return nil, fmt.Errorf("%s: %s", errorSummary, apiResp.Message.String())
	}

	return apiResp, nil
}

// Polling constants: balance fewer GET calls with quick result delivery.
// - initialDelay: wait before first GET (backend rarely has result in first few hundred ms; avoids a GET that almost always returns 202).
// - initialInterval: first wait between polls after a 202.
// - backoffMultiplier / maxInterval: exponential backoff to limit GET rate for long-running commands.
const (
	pollInitialDelay      = 400 * time.Millisecond
	pollInitialInterval   = 800 * time.Millisecond
	pollBackoffMultiplier = 1.5
	pollMaxInterval       = 4 * time.Second
)

// PollRequestResult polls RequestResult until the request is completed.
// Uses an initial delay (skip immediate GET that usually returns 202), then exponential backoff.
func (a *ApiManager) PollRequestResult(ctx context.Context, requestId string, timeout ...time.Duration) (*RequestResultResponse, error) {
	pollTimeout := a.RequestTimeout
	if pollTimeout <= 0 {
		pollTimeout = 60 * time.Second
	}
	if len(timeout) > 0 && timeout[0] > 0 {
		pollTimeout = timeout[0]
	}

	pollCtx, cancel := context.WithTimeout(ctx, pollTimeout)
	defer cancel()

	currentInterval := pollInitialInterval

	// Wait briefly before first GET (most commands are not ready in the first few hundred ms)
	select {
	case <-pollCtx.Done():
		return nil, fmt.Errorf("timeout waiting for request %s to complete (timeout: %v)", requestId, pollTimeout)
	case <-time.After(pollInitialDelay):
	}

	for {
		result, err := a.RequestResult(pollCtx, requestId)
		if err == nil && result.Status == "success" {
			return result, nil
		}
		if err != nil && !strings.Contains(err.Error(), "still in queued/processing") {
			return nil, err
		}

		select {
		case <-pollCtx.Done():
			return nil, fmt.Errorf("timeout waiting for request %s to complete (timeout: %v)", requestId, pollTimeout)
		case <-time.After(currentInterval):
		}

		currentInterval = min(time.Duration(float64(currentInterval)*pollBackoffMultiplier), pollMaxInterval)
	}
}

// DetectAccountType detects whether the account is MSP or Enterprise.
func (a *ApiManager) IsMspAccountType(ctx context.Context) error {
	// Try a lightweight MSP command to detect account type
	// If it succeeds, it's an MSP account; if it fails, it's Enterprise
	apiResp, err := a.ExecuteCommand(ctx, "switch-to-msp", "Account type detection")

	if err != nil {
		// Check if error is specifically about MSP not being available
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "already msp") ||
			strings.Contains(errMsg, "already") {
			a.IsMspAccount = true
			return nil
		}

		if strings.Contains(errMsg, "restricted") {
			a.IsMspAccount = false
		}

		return nil
	}

	if strings.Contains(apiResp.Message.String(), "restricted") {
		a.IsMspAccount = false
	}

	// Command succeeded - it's an MSP account
	a.IsMspAccount = true
	return nil
}
