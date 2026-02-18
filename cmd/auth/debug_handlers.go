package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sapliy/fintech-ecosystem/internal/flow"
	flowDomain "github.com/sapliy/fintech-ecosystem/internal/flow/domain"
	"github.com/sapliy/fintech-ecosystem/pkg/jsonutil"
)

type DebugHandler struct {
	debugService *flow.DebugService
}

func NewDebugHandler(debugService *flow.DebugService) *DebugHandler {
	return &DebugHandler{debugService: debugService}
}

// StartDebugSession starts a new debug session
func (h *DebugHandler) StartDebugSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FlowID string `json:"flow_id"`
		ZoneID string `json:"zone_id"`
		Level  string `json:"level"` // info, verbose, trace
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.WriteErrorJSON(w, "Invalid request body")
		return
	}

	// Validate debug level
	level := flowDomain.DebugLevel(req.Level)
	if level != flowDomain.DebugLevelInfo && level != flowDomain.DebugLevelVerbose && level != flowDomain.DebugLevelTrace {
		level = flowDomain.DebugLevelInfo
	}

	session, err := h.debugService.StartDebugSession(context.Background(), req.FlowID, req.ZoneID, level)
	if err != nil {
		jsonutil.WriteErrorJSON(w, err.Error())
		return
	}

	jsonutil.WriteJSON(w, http.StatusCreated, session)
}

// GetDebugSession retrieves a debug session
func (h *DebugHandler) GetDebugSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		jsonutil.WriteErrorJSON(w, "session_id query parameter required")
		return
	}

	session, err := h.debugService.GetDebugSession(sessionID)
	if err != nil {
		jsonutil.WriteErrorJSON(w, err.Error())
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, session)
}

// GetDebugEvents retrieves debug events for a session
func (h *DebugHandler) GetDebugEvents(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		jsonutil.WriteErrorJSON(w, "session_id query parameter required")
		return
	}

	// Parse optional since parameter
	var since *time.Time
	if sinceStr := r.URL.Query().Get("since"); sinceStr != "" {
		if timestamp, err := strconv.ParseInt(sinceStr, 10, 64); err == nil {
			t := time.Unix(timestamp/1000, 0) // Convert from milliseconds if needed
			since = &t
		}
	}

	events, err := h.debugService.GetDebugEvents(sessionID, since)
	if err != nil {
		jsonutil.WriteErrorJSON(w, err.Error())
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, events)
}

// EndDebugSession ends a debug session
func (h *DebugHandler) EndDebugSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID string `json:"session_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.WriteErrorJSON(w, "Invalid request body")
		return
	}

	if err := h.debugService.EndDebugSession(req.SessionID); err != nil {
		jsonutil.WriteErrorJSON(w, err.Error())
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ended",
	})
}

// WebSocketDebug handles WebSocket connections for real-time debug events
func (h *DebugHandler) WebSocketDebug(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "session_id query parameter required", http.StatusBadRequest)
		return
	}

	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade to WebSocket", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Get debug session
	_, err = h.debugService.GetDebugSession(sessionID)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"debug session not found"}`))
		return
	}

	// Send existing events since last check (for reconnection)
	events, err := h.debugService.GetDebugEvents(sessionID, nil)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"could not retrieve events"}`))
		return
	}

	// Send existing events
	for _, event := range events {
		eventJSON, _ := json.Marshal(event)
		if err := conn.WriteMessage(websocket.TextMessage, eventJSON); err != nil {
			return
		}
	}

	// Keep connection alive and stream new events
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Check for new events
		oneHourAgo := time.Now().Add(-1 * time.Hour)
		newEvents, err := h.debugService.GetDebugEvents(sessionID, &oneHourAgo)
		if err != nil {
			continue
		}

		for _, event := range newEvents {
			eventJSON, _ := json.Marshal(event)
			if err := conn.WriteMessage(websocket.TextMessage, eventJSON); err != nil {
				return
			}
		}
	}
}

// ExecuteFlowWithDebug executes a flow in debug mode
func (h *DebugHandler) ExecuteFlowWithDebug(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FlowID       string                 `json:"flow_id"`
		DebugSession string                 `json:"debug_session_id"`
		Input        map[string]interface{} `json:"input"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.WriteErrorJSON(w, "Invalid request body")
		return
	}

	// Get flow and debug session
	ctx := context.Background()

	// Create debug runner
	repo := h.debugService.GetRepository()
	baseRunner := flowDomain.NewFlowRunner(repo)
	debugRunner := flow.NewDebugFlowRunner(baseRunner, h.debugService, repo)

	// Get flow
	flow, err := repo.GetFlow(ctx, req.FlowID)
	if err != nil {
		jsonutil.WriteErrorJSON(w, "Flow not found")
		return
	}

	// Execute with debug
	err = debugRunner.ExecuteWithDebug(ctx, flow, req.Input, req.DebugSession)
	if err != nil {
		jsonutil.WriteErrorJSON(w, err.Error())
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "debug_execution_started",
	})
}
