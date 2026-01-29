package fraud

import (
	"context"
	"fmt"
)

type RuleResult struct {
	RuleName string
	Passed   bool
	Message  string
}

type Rule interface {
	Name() string
	Check(ctx context.Context, tx Transaction) (RuleResult, error)
}

type Transaction struct {
	ID       string
	Amount   int64
	Currency string
	UserID   string
}

type Engine struct {
	rules []Rule
}

func NewEngine(rules ...Rule) *Engine {
	return &Engine{rules: rules}
}

func (e *Engine) Check(ctx context.Context, tx Transaction) ([]RuleResult, bool) {
	results := make([]RuleResult, 0, len(e.rules))
	isRisky := false

	for _, rule := range e.rules {
		res, err := rule.Check(ctx, tx)
		if err != nil {
			results = append(results, RuleResult{
				RuleName: rule.Name(),
				Passed:   false,
				Message:  fmt.Sprintf("Error: %v", err),
			})
			continue
		}
		results = append(results, res)
		if !res.Passed {
			isRisky = true
		}
	}

	return results, isRisky
}
