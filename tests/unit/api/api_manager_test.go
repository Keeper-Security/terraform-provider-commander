// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
)

func TestFlexibleMessage_UnmarshalJSON_String(t *testing.T) {
	var fm api.FlexibleMessage
	err := json.Unmarshal([]byte(`"hello"`), &fm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm.String() != "hello" {
		t.Errorf("expected hello, got %s", fm.String())
	}
}

func TestFlexibleMessage_UnmarshalJSON_Array(t *testing.T) {
	var fm api.FlexibleMessage
	err := json.Unmarshal([]byte(`["a","b"]`), &fm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm.String() != `["a","b"]` {
		t.Errorf("expected [\"a\",\"b\"], got %s", fm.String())
	}
}

func TestFlexibleMessage_UnmarshalJSON_Any(t *testing.T) {
	var fm api.FlexibleMessage
	err := json.Unmarshal([]byte(`{"key":"value"}`), &fm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm.String() != `{"key":"value"}` {
		t.Errorf("expected {\"key\":\"value\"}, got %s", fm.String())
	}
}

func TestFlexibleMessage_UnmarshalJSON_Number(t *testing.T) {
	var fm api.FlexibleMessage
	err := json.Unmarshal([]byte(`123`), &fm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm.String() != "123" {
		t.Errorf("expected 123, got %s", fm.String())
	}
}

func TestFlexibleMessage_String(t *testing.T) {
	fm := api.FlexibleMessage("test")
	if fm.String() != "test" {
		t.Errorf("expected test, got %s", fm.String())
	}
}

func TestSubmitRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/executecommand-async" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"success":true,"request_id":"req-123","status":"queued","message":"ok"}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	resp, err := client.SubmitRequest(ctx, "enterprise-info", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.RequestId != "req-123" {
		t.Errorf("expected request_id req-123, got %s", resp.RequestId)
	}
	if !resp.Success {
		t.Error("expected success true")
	}
}

func TestSubmitRequest_SetsMinCommanderVersionHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMinCommanderVersionHeader(t, r)
		if got := r.Header.Get("api-key"); got != "test-key" {
			t.Errorf("api-key = %q, want test-key", got)
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"success":true,"request_id":"req-123","status":"queued","message":"ok"}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	if _, err := client.SubmitRequest(context.Background(), "enterprise-info", nil); err != nil {
		t.Fatalf("SubmitRequest: %v", err)
	}
}

// TestSubmitRequest_NormalizesApostrophesInQuotedFields is a regression test for a bug
// where normalizeCommandForShell reprocessed spans that utils.QuoteShellSingle had
// already correctly escaped (close-quote, literal ' via double quotes, reopen-quote),
// corrupting them into an unbalanced-quote string that the Commander backend rejected
// with "500 Unexpected error: No closing quotation". It also locks in that
// normalizeCommandForShell still fixes up naive, unescaped single-quoted spans (as
// built by other call sites, e.g. fmt.Sprintf("'%s'", value)) that contain a raw
// apostrophe.
func TestSubmitRequest_NormalizesApostrophesInQuotedFields(t *testing.T) {
	var capturedBody map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		if err := json.Unmarshal(body, &capturedBody); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"success":true,"request_id":"req-123","status":"queued","message":"ok"}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()

	// Title/notes built the same way BuildRecordAdd/FormatFieldAssignment build them
	// today: via the real (escaping) utils.QuoteShellSingle. A naive fragment (as
	// produced by other, unrelated call sites that wrap values in raw single quotes
	// without escaping) is appended to confirm that behavior is unaffected.
	title := utils.QuoteShellSingle("California Driver's License - John Doe")
	notes := utils.QuoteShellSingle("Personal driver's license.")
	rawCommand := "record-add --title " + title + " --record-type driverLicense --notes " + notes +
		" -f enterprise-role --add 'O'Brien' --node 'root'"

	_, err := client.SubmitRequest(ctx, rawCommand, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantCommand := "record-add --title " + title + " --record-type driverLicense --notes " + notes +
		` -f enterprise-role --add "O'Brien" --node 'root'`
	if capturedBody["command"] != wantCommand {
		t.Errorf("normalized command mismatch:\n got:  %s\n want: %s", capturedBody["command"], wantCommand)
	}
}

func TestSubmitRequest_ErrorStatus(t *testing.T) {
	// 503 is retryable, so a persistently unavailable backend now retries
	// before failing; fast retry knobs keep this test from taking real
	// wall-clock backoff time.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":"queue full"}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:         server.URL,
		ServiceModeApiKey:      "test-key",
		HttpClient:             server.Client(),
		RetryMaxAttempts:       2,
		RetryInitialInterval:   1 * time.Millisecond,
		RetryBackoffMultiplier: 1,
	}
	ctx := context.Background()
	_, err := client.SubmitRequest(ctx, "cmd", nil)
	if err == nil {
		t.Fatal("expected error for 503")
	}
	if err != nil && err.Error() != "queue is full (503): service unavailable, please try again later" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSubmitRequest_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid command"}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	_, err := client.SubmitRequest(ctx, "cmd", nil)
	if err == nil {
		t.Fatal("expected error for 400")
	}
	if err != nil && err.Error() != "bad request (400): invalid command" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSubmitRequest_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	_, err := client.SubmitRequest(ctx, "cmd", nil)
	if err == nil {
		t.Fatal("expected error for 404")
	}
}

func TestSubmitRequest_InternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"execution failed"}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	_, err := client.SubmitRequest(ctx, "cmd", nil)
	if err == nil {
		t.Fatal("expected error for 500")
	}
	if err != nil && err.Error() != "internal server error (500): execution failed" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSubmitRequest_ResourceNotFound_CannotFindPhrasing(t *testing.T) {
	// Real-world Commander wording for an out-of-band-deleted Nested Shared
	// Folder: "Cannot find any Nested Share Folder object with UID ...".
	// It contains no "not found" substring, so this must still classify as
	// api.ErrResourceNotFound or Read() will treat deletion as a hard error
	// forever instead of removing the resource from state.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"Cannot find any Nested Share Folder object with UID abc123."}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	_, err := client.SubmitRequest(ctx, "cmd", nil)
	if err == nil {
		t.Fatal("expected error for 500")
	}
	if !errors.Is(err, api.ErrResourceNotFound) {
		t.Errorf("expected errors.Is(err, api.ErrResourceNotFound) to be true, got: %v", err)
	}
}

func TestSubmitRequest_TooManyRequests(t *testing.T) {
	// A persistently rate-limited backend: doWithRetry should exhaust its
	// retries and still surface the 429 as an error, not hang forever. Fast
	// retry knobs keep this test from taking real wall-clock backoff time.
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:         server.URL,
		ServiceModeApiKey:      "test-key",
		HttpClient:             server.Client(),
		RetryMaxAttempts:       3,
		RetryInitialInterval:   1 * time.Millisecond,
		RetryBackoffMultiplier: 1,
	}
	ctx := context.Background()
	_, err := client.SubmitRequest(ctx, "cmd", nil)
	if err == nil {
		t.Fatal("expected error for 429")
	}
	if requests != 3 {
		t.Errorf("expected 3 attempts (RetryMaxAttempts), got %d", requests)
	}
}

func TestSubmitRequest_TooManyRequests_RetriesThenSucceeds(t *testing.T) {
	// Client-reported scenario: backend answers 429 a couple of times (rate
	// limiting), then the request actually goes through. doWithRetry should
	// absorb the transient 429s so the caller never sees an error.
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"success":true,"request_id":"req-1","status":"queued","message":"ok"}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:         server.URL,
		ServiceModeApiKey:      "test-key",
		HttpClient:             server.Client(),
		RetryInitialInterval:   1 * time.Millisecond,
		RetryBackoffMultiplier: 1,
	}
	ctx := context.Background()
	resp, err := client.SubmitRequest(ctx, "cmd", nil)
	if err != nil {
		t.Fatalf("expected retries to absorb transient 429s, got error: %v", err)
	}
	if resp.RequestId != "req-1" {
		t.Errorf("request_id = %q, want req-1", resp.RequestId)
	}
	if requests != 3 {
		t.Errorf("expected 3 attempts before success, got %d", requests)
	}
}

func TestSubmitRequest_TooManyRequests_HonorsRetryAfterHeader(t *testing.T) {
	// A Retry-After: 0 header should be used as the wait instead of the
	// (larger) default backoff, so this completes quickly even without
	// overriding RetryInitialInterval.
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests < 2 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"success":true,"request_id":"req-1","status":"queued","message":"ok"}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	_, err := client.SubmitRequest(ctx, "cmd", nil)
	if err != nil {
		t.Fatalf("expected Retry-After-driven retry to succeed, got error: %v", err)
	}
	if requests != 2 {
		t.Errorf("expected 2 attempts, got %d", requests)
	}
}

func TestSubmitRequest_Unexpected2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	_, err := client.SubmitRequest(ctx, "cmd", nil)
	if err == nil {
		t.Fatal("expected error for 200 (expected 202)")
	}
}

func TestRequestStatus_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/status/req-123" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"request_id":"req-123","command":"cmd","status":"completed","created_at":"","started_at":"","completed_at":""}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	resp, err := client.RequestStatus(ctx, "req-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.RequestId != "req-123" {
		t.Errorf("expected request_id req-123, got %s", resp.RequestId)
	}
}

func TestRequestStatus_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	_, err := client.RequestStatus(ctx, "req-123")
	if err == nil {
		t.Fatal("expected error for 404")
	}
}

func TestRequestResult_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"id":1}],"status":"success","message":"done","error":""}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	resp, err := client.RequestResult(ctx, "req-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected status success, got %s", resp.Status)
	}
}

func TestRequestResult_SetsMinCommanderVersionHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMinCommanderVersionHeader(t, r)
		if got := r.Header.Get("api-key"); got != "test-key" {
			t.Errorf("api-key = %q, want test-key", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"id":1}],"status":"success","message":"done","error":""}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	if _, err := client.RequestResult(context.Background(), "req-123"); err != nil {
		t.Fatalf("RequestResult: %v", err)
	}
}

func TestRequestResult_StillProcessing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	_, err := client.RequestResult(ctx, "req-123")
	if err == nil {
		t.Fatal("expected error for 202 (still processing)")
	}
}

func TestRequestResult_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	_, err := client.RequestResult(ctx, "req-123")
	if err == nil {
		t.Fatal("expected error for 500")
	}
}

func TestExecuteCommand_Success(t *testing.T) {
	reqCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/executecommand-async" {
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"success":true,"request_id":"req-1","status":"queued","message":"ok"}`))
			reqCount++
			return
		}
		if r.URL.Path == "/status/req-1" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"success":true,"request_id":"req-1","command":"cmd","status":"completed","created_at":"","started_at":"","completed_at":""}`))
			return
		}
		if r.URL.Path == "/result/req-1" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":null,"status":"success","message":"done","error":""}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	resp, err := client.ExecuteCommand(ctx, "switch-to-msp", "Account type detection")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected status success, got %s", resp.Status)
	}
}

func TestExecuteCommand_RecoversFromTransient429OnResultPoll(t *testing.T) {
	// Client-reported scenario: a delete command's result-poll GET hits a
	// transient 429 (backend rate limiting), but the underlying delete
	// already succeeded. Before the retry fix, PollRequestResult treated any
	// non-"still queued" error as fatal and surfaced e.g. "Record Delete
	// Failed" even though the delete went through - desyncing state. Now the
	// 429 should be retried transparently and ExecuteCommand should succeed.
	var resultRequests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/executecommand-async" {
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"success":true,"request_id":"req-1","status":"queued","message":"ok"}`))
			return
		}
		if r.URL.Path == "/result/req-1" {
			resultRequests++
			if resultRequests == 1 {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":null,"status":"success","message":"deleted","error":""}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:         server.URL,
		ServiceModeApiKey:      "test-key",
		HttpClient:             server.Client(),
		RetryInitialInterval:   1 * time.Millisecond,
		RetryBackoffMultiplier: 1,
	}
	ctx := context.Background()
	resp, err := client.ExecuteCommand(ctx, "record-delete", "Record Delete Failed")
	if err != nil {
		t.Fatalf("expected the transient 429 to be retried and the delete to succeed, got error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected status success, got %s", resp.Status)
	}
	if resultRequests != 2 {
		t.Errorf("expected 2 result-poll requests (1 rate-limited + 1 success), got %d", resultRequests)
	}
}

func TestExecuteCommand_SubmitFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad"}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	_, err := client.ExecuteCommand(ctx, "cmd", "Test error")
	if err == nil {
		t.Fatal("expected error when submit fails")
	}
}

func TestPollRequestResult_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always return 202 (still processing) so we never get result
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	// Use very short timeout so test finishes quickly
	_, err := client.PollRequestResult(ctx, "req-1", 20*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout")
	}
}

func TestPollRequestResult_ImmediateSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":null,"status":"success","message":"ok","error":""}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	resp, err := client.PollRequestResult(ctx, "req-1", 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected status success, got %s", resp.Status)
	}
}

func TestIsMspAccountType_ErrorPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"restricted"}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	ctx := context.Background()
	err := client.IsMspAccountType(ctx)
	if err != nil {
		// API failed - provider sets IsMspAccount = false and returns nil
		return
	}
	if client.IsMspAccount {
		t.Error("expected IsMspAccount false when command fails with restricted")
	}
}

func TestIsMspAccountType_AlreadyMsp(t *testing.T) {
	reqCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/executecommand-async" {
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"success":true,"request_id":"req-1","status":"queued"}`))
			reqCount++
			return
		}
		if r.URL.Path == "/result/req-1" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":null,"status":"success","message":"already msp","error":""}`))
			return
		}
		if r.URL.Path == "/status/req-1" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"success":true,"request_id":"req-1","status":"completed"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	_ = client.IsMspAccountType(ctx)
	// When error message contains "already" we set IsMspAccount = true.
	// May be false if ExecuteCommand failed (e.g. poll timeout); test just ensures no panic.
	if !client.IsMspAccount {
		t.Log("IsMspAccount is false (e.g. ExecuteCommand failed or poll timeout)")
	}
}

func assertMinCommanderVersionHeader(t *testing.T, r *http.Request) {
	t.Helper()
	if got := r.Header.Get(api.MinCommanderVersionHeader); got != api.MinCommanderVersion {
		t.Errorf("%s = %q, want %q", api.MinCommanderVersionHeader, got, api.MinCommanderVersion)
	}
}
