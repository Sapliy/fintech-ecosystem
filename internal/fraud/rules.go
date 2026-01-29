package fraud

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// AmountRule checks if a transaction amount exceeds a limit.
type AmountRule struct {
	Limit int64
}

func (r *AmountRule) Name() string { return "AmountRule" }

func (r *AmountRule) Check(ctx context.Context, tx Transaction) (RuleResult, error) {
	if tx.Amount > r.Limit {
		return RuleResult{
			RuleName: r.Name(),
			Passed:   false,
			Message:  fmt.Sprintf("Amount %d exceeds limit %d", tx.Amount, r.Limit),
		}, nil
	}
	return RuleResult{RuleName: r.Name(), Passed: true}, nil
}

// VelocityRule checks if a user has made too many transactions in a time window.
type VelocityRule struct {
	Window    time.Duration
	Threshold int
	mu        sync.Mutex
	history   map[string][]time.Time
}

func NewVelocityRule(window time.Duration, threshold int) *VelocityRule {
	return &VelocityRule{
		Window:    window,
		Threshold: threshold,
		history:   make(map[string][]time.Time),
	}
}

func (r *VelocityRule) Name() string { return "VelocityRule" }

func (r *VelocityRule) Check(ctx context.Context, tx Transaction) (RuleResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	r.history[tx.UserID] = append(r.history[tx.UserID], now)

	var fresh []time.Time
	for _, t := range r.history[tx.UserID] {
		if now.Sub(t) < r.Window {
			fresh = append(fresh, t)
		}
	}
	r.history[tx.UserID] = fresh

	if len(fresh) > r.Threshold {
		return RuleResult{
			RuleName: r.Name(),
			Passed:   false,
			Message:  fmt.Sprintf("Velocity high: %d transactions in %v", len(fresh), r.Window),
		}, nil
	}

	return RuleResult{RuleName: r.Name(), Passed: true}, nil
}
