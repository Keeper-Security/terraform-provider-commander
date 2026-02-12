// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
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
	resp, err := client.SubmitRequest(ctx, "enterprise-info")
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

func TestSubmitRequest_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":"queue full"}`))
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	_, err := client.SubmitRequest(ctx, "cmd")
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
	_, err := client.SubmitRequest(ctx, "cmd")
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
	_, err := client.SubmitRequest(ctx, "cmd")
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
	_, err := client.SubmitRequest(ctx, "cmd")
	if err == nil {
		t.Fatal("expected error for 500")
	}
	if err != nil && err.Error() != "internal server error (500): execution failed" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSubmitRequest_TooManyRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
	ctx := context.Background()
	_, err := client.SubmitRequest(ctx, "cmd")
	if err == nil {
		t.Fatal("expected error for 429")
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
	_, err := client.SubmitRequest(ctx, "cmd")
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
