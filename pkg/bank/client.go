package bank

import (
	"context"
	"errors"
	"time"
)

type Status string

const (
	StatusSuccess Status = "succeeded"
	StatusFailed  Status = "failed"
)

// TransactionResult represents the outcome of a bank transaction.
type TransactionResult struct {
	TransactionID string
	Status        Status
	ErrorCode     string
}

// Client defines the interface for communicating with a bank.
type Client interface {
	Charge(ctx context.Context, amount int64, currency, cardToken string) (*TransactionResult, error)
}

// MockClient is a mock implementation of the Bank Client.
type MockClient struct{}

// NewMockClient creates a new instance of MockClient.
func NewMockClient() *MockClient {
	return &MockClient{}
}

// Charge simulates a credit card charge.
// Special logic for testing:
// - specific amounts can trigger failures.
// - cardToken "tok_visa" -> success
// - cardToken "tok_mastercard" -> success
// - other tokens -> failure
func (m *MockClient) Charge(ctx context.Context, amount int64, currency, cardToken string) (*TransactionResult, error) {
	// Simulate network latency
	time.Sleep(500 * time.Millisecond)

	// Validation
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}

	// Simulation logic
	switch cardToken {
	case "tok_visa", "tok_mastercard":
		return &TransactionResult{
			TransactionID: "txn_" + GenerateRandomID(),
			Status:        "succeeded",
		}, nil
	case "tok_declined":
		return &TransactionResult{
			Status:    "failed",
			ErrorCode: "card_declined",
		}, nil
	default:
		// Default to success for development ease, or error?
		// Let's default to error to be safe and force correct token usage.
		return &TransactionResult{
			Status:    "failed",
			ErrorCode: "invalid_card_token",
		}, nil
	}
}

func GenerateRandomID() string {
	// In real app use crypto/rand or uuid
	return time.Now().Format("20060102150405")
}
