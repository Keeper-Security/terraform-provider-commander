// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
)

// CommandServer queues commands and returns request_id req-1, req-2, ... for GET /result.
// Used by unit tests to mock the Commander executecommand-async API.
type CommandServer struct {
	mu       sync.Mutex
	commands []string
}

// CommandCount returns the number of commands received (POST executecommand-async calls).
func (c *CommandServer) CommandCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.commands)
}

// StartCommandServer starts an httptest.Server that handles the ExecuteCommand flow:
// POST .../executecommand-async (body: {"command":"..."}) -> 202 + request_id;
// GET .../result/req-N -> 200 + JSON with message and data from responseForCommand.
// responseForCommand(cmd, idx) is 1-based; pass nil to use default "ok", nil.
func StartCommandServer(mock *CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	if responseForCommand == nil {
		responseForCommand = func(string, int) (string, interface{}) { return "ok", nil }
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/executecommand-async") {
			var body struct {
				Command string `json:"command"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			mock.mu.Lock()
			mock.commands = append(mock.commands, body.Command)
			n := len(mock.commands)
			mock.mu.Unlock()
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"success":true,"request_id":"req-` + strconv.Itoa(n) + `","status":"queued"}`))
			return
		}
		if r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/result/") {
			parts := strings.Split(r.URL.Path, "/result/")
			idx := 1
			if len(parts) == 2 {
				if n, err := strconv.Atoi(strings.TrimPrefix(parts[1], "req-")); err == nil {
					idx = n
				}
			}
			mock.mu.Lock()
			cmd := "ok"
			if idx <= len(mock.commands) {
				cmd = mock.commands[idx-1]
			}
			mock.mu.Unlock()
			msg, data := responseForCommand(cmd, idx)
			dataBytes, _ := json.Marshal(data)
			if dataBytes == nil {
				dataBytes = []byte("null")
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":` + string(dataBytes) + `,"status":"success","message":"` + msg + `","error":""}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}
