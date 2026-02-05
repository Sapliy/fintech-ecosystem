package domain

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// DebugLevel represents the level of debug information
type DebugLevel string

const (
	DebugLevelInfo    DebugLevel = "info"
	DebugLevelVerbose DebugLevel = "verbose"
	DebugLevelTrace   DebugLevel = "trace"
)

// DebugEvent represents a debug event during flow execution
type DebugEvent struct {
	ID          string                 `json:"id"`
	ExecutionID string                 `json:"execution_id"`
	FlowID      string                 `json:"flow_id"`
	NodeID      string                 `json:"node_id"`
	Level       DebugLevel             `json:"level"`
	Type        DebugEventType         `json:"type"`
	Message     string                 `json:"message"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// DebugEventType represents types of debug events
type DebugEventType string

const (
	DebugEventExecutionStart DebugEventType = "execution_start"
	DebugEventExecutionEnd   DebugEventType = "execution_end"
	DebugEventNodeStart      DebugEventType = "node_start"
	DebugEventNodeEnd        DebugEventType = "node_end"
	DebugEventNodeError      DebugEventType = "node_error"
	DebugEventNodePaused     DebugEventType = "node_paused"
	DebugEventConditionEval  DebugEventType = "condition_eval"
	DebugEventWebhookCall    DebugEventType = "webhook_call"
	DebugEventApprovalReq    DebugEventType = "approval_required"
)

// DebugSession represents an active debug session
type DebugSession struct {
	ID        string            `json:"id"`
	FlowID    string            `json:"flow_id"`
	ZoneID    string            `json:"zone_id"`
	Level     DebugLevel        `json:"level"`
	Active    bool              `json:"active"`
	StartTime time.Time         `json:"start_time"`
	Events    []DebugEvent      `json:"events"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
}

// DebugSessionManager manages active debug sessions
type DebugSessionManager struct {
	sessions map[string]*DebugSession
	mu       sync.RWMutex
}

// NewDebugSessionManager creates a new debug session manager
func NewDebugSessionManager() *DebugSessionManager {
	return &DebugSessionManager{
		sessions: make(map[string]*DebugSession),
	}
}

// CreateSession creates a new debug session
func (m *DebugSessionManager) CreateSession(ctx context.Context, flowID, zoneID string, level DebugLevel) (*DebugSession, error) {
	sessionID := fmt.Sprintf("debug_%d", time.Now().UnixNano())

	session := &DebugSession{
		ID:        sessionID,
		FlowID:    flowID,
		ZoneID:    zoneID,
		Level:     level,
		Active:    true,
		StartTime: time.Now(),
		Events:    make([]DebugEvent, 0),
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
	}

	m.mu.Lock()
	m.sessions[sessionID] = session
	m.mu.Unlock()

	// Send start event
	m.logEvent(sessionID, DebugEventExecutionStart, DebugLevelInfo, "Flow execution started", map[string]interface{}{
		"flow_id": flowID,
		"zone_id": zoneID,
	})

	return session, nil
}

// GetSession retrieves a debug session
func (m *DebugSessionManager) GetSession(sessionID string) (*DebugSession, error) {
	m.mu.RLock()
	session, exists := m.sessions[sessionID]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("debug session not found: %s", sessionID)
	}

	return session, nil
}

// EndSession ends a debug session
func (m *DebugSessionManager) EndSession(sessionID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[sessionID]
	if !exists {
		return fmt.Errorf("debug session not found: %s", sessionID)
	}

	session.Active = false

	// Send end event
	m.logEventUnsafe(sessionID, DebugEventExecutionEnd, DebugLevelInfo, "Flow execution ended", map[string]interface{}{
		"duration": time.Since(session.StartTime).String(),
	})

	return nil
}

// LogNodeStart logs when a node starts execution
func (m *DebugSessionManager) LogNodeStart(sessionID, nodeID, nodeType string, input map[string]interface{}) {
	data := map[string]interface{}{
		"node_id":   nodeID,
		"node_type": nodeType,
		"input":     input,
	}
	m.logEvent(sessionID, DebugEventNodeStart, DebugLevelInfo, fmt.Sprintf("Starting execution of node %s", nodeID), data)
}

// LogNodeEnd logs when a node completes execution
func (m *DebugSessionManager) LogNodeEnd(sessionID, nodeID string, output map[string]interface{}, duration time.Duration) {
	data := map[string]interface{}{
		"node_id":  nodeID,
		"output":   output,
		"duration": duration.String(),
	}
	m.logEvent(sessionID, DebugEventNodeEnd, DebugLevelInfo, fmt.Sprintf("Completed execution of node %s", nodeID), data)
}

// LogNodeError logs when a node encounters an error
func (m *DebugSessionManager) LogNodeError(sessionID, nodeID string, err error) {
	data := map[string]interface{}{
		"node_id": nodeID,
		"error":   err.Error(),
	}
	m.logEvent(sessionID, DebugEventNodeError, DebugLevelInfo, fmt.Sprintf("Node %s encountered error: %s", nodeID, err.Error()), data)
}

// LogNodePaused logs when a node is paused
func (m *DebugSessionManager) LogNodePaused(sessionID, nodeID string, reason string) {
	data := map[string]interface{}{
		"node_id": nodeID,
		"reason":  reason,
	}
	m.logEvent(sessionID, DebugEventNodePaused, DebugLevelInfo, fmt.Sprintf("Node %s paused: %s", nodeID, reason), data)
}

// LogConditionEval logs condition evaluation
func (m *DebugSessionManager) LogConditionEval(sessionID, nodeID, field, operator string, value, inputValue interface{}, result bool) {
	data := map[string]interface{}{
		"node_id":     nodeID,
		"field":       field,
		"operator":    operator,
		"value":       value,
		"input_value": inputValue,
		"result":      result,
	}
	level := DebugLevelInfo
	if result {
		level = DebugLevelVerbose
	}
	m.logEvent(sessionID, DebugEventConditionEval, level, fmt.Sprintf("Condition evaluated: %s %s %v = %v", field, operator, value, result), data)
}

// LogWebhookCall logs webhook execution
func (m *DebugSessionManager) LogWebhookCall(sessionID, nodeID, url, method string, payload map[string]interface{}) {
	data := map[string]interface{}{
		"node_id": nodeID,
		"url":     url,
		"method":  method,
		"payload": payload,
	}
	m.logEvent(sessionID, DebugEventWebhookCall, DebugLevelVerbose, fmt.Sprintf("Webhook called: %s %s", method, url), data)
}

// LogApprovalRequired logs when approval is required
func (m *DebugSessionManager) LogApprovalRequired(sessionID, nodeID, approver string) {
	data := map[string]interface{}{
		"node_id":  nodeID,
		"approver": approver,
	}
	m.logEvent(sessionID, DebugEventApprovalReq, DebugLevelInfo, fmt.Sprintf("Approval required for node %s by %s", nodeID, approver), data)
}

// GetEvents returns events for a session, optionally filtered by level
func (m *DebugSessionManager) GetEvents(sessionID string, since *time.Time) ([]DebugEvent, error) {
	m.mu.RLock()
	session, exists := m.sessions[sessionID]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("debug session not found: %s", sessionID)
	}

	var events []DebugEvent
	for _, event := range session.Events {
		if since != nil && event.Timestamp.Before(*since) {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

// GetActiveSessions returns all active debug sessions for a flow
func (m *DebugSessionManager) GetActiveSessions(flowID string) []*DebugSession {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var sessions []*DebugSession
	for _, session := range m.sessions {
		if session.Active && session.FlowID == flowID {
			sessions = append(sessions, session)
		}
	}

	return sessions
}

// CleanupOldSessions removes sessions older than the specified duration
func (m *DebugSessionManager) CleanupOldSessions(maxAge time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for id, session := range m.sessions {
		if now.Sub(session.CreatedAt) > maxAge {
			delete(m.sessions, id)
			log.Printf("Cleaned up old debug session: %s", id)
		}
	}
}

// Private helper methods

func (m *DebugSessionManager) logEvent(sessionID string, eventType DebugEventType, level DebugLevel, message string, data map[string]interface{}) {
	eventID := fmt.Sprintf("event_%d", time.Now().UnixNano())

	event := DebugEvent{
		ID:          eventID,
		ExecutionID: sessionID,
		Type:        eventType,
		Level:       level,
		Message:     message,
		Data:        data,
		Timestamp:   time.Now(),
	}

	// Extract NodeID from data for node-related events
	if data != nil {
		if nodeID, ok := data["node_id"]; ok {
			if nodeIDStr, ok := nodeID.(string); ok {
				event.NodeID = nodeIDStr
			}
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.sessions[sessionID]; exists {
		session.Events = append(session.Events, event)
	} else {
		// Session doesn't exist, create a temporary one for logging
		log.Printf("Warning: Debug session %s not found for event: %s", sessionID, message)
	}
}

func (m *DebugSessionManager) logEventUnsafe(sessionID string, eventType DebugEventType, level DebugLevel, message string, data map[string]interface{}) {
	eventID := fmt.Sprintf("event_%d", time.Now().UnixNano())

	event := DebugEvent{
		ID:          eventID,
		ExecutionID: sessionID,
		Type:        eventType,
		Level:       level,
		Message:     message,
		Data:        data,
		Timestamp:   time.Now(),
	}

	if session, exists := m.sessions[sessionID]; exists {
		session.Events = append(session.Events, event)
	}
}

// ShouldLog determines if an event should be logged based on session level
func (m *DebugSessionManager) ShouldLog(sessionID string, eventLevel DebugLevel) bool {
	m.mu.RLock()
	session, exists := m.sessions[sessionID]
	m.mu.RUnlock()

	if !exists {
		return false
	}

	// All levels are logged, but filtered on retrieval
	return session.Active
}
