package policy

import (
	"context"
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/rego"
)

// OPAPolicyEngine implements PolicyEngine using Open Policy Agent (Rego)
type OPAPolicyEngine struct {
	query rego.PreparedEvalQuery
}

// NewOPAPolicyEngine creates a new OPA policy engine from a Rego file
func NewOPAPolicyEngine(path string) (*OPAPolicyEngine, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read rego file: %w", err)
	}

	ctx := context.Background()
	query, err := rego.New(
		rego.Query("data.sapliy.authz.allow"),
		rego.Module(path, string(data)),
	).PrepareForEval(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to prepare rego query: %w", err)
	}

	return &OPAPolicyEngine{query: query}, nil
}

// Check evaluates policies using OPA Rego
func (e *OPAPolicyEngine) Check(ctx context.Context, pctx *PolicyContext) (*PolicyResult, error) {
	// Prepare input for OPA
	input := map[string]interface{}{
		"roles":  pctx.Roles,
		"action": string(pctx.Action),
	}

	results, err := e.query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate policy: %w", err)
	}

	allowed := false
	if len(results) > 0 {
		if val, ok := results[0].Expressions[0].Value.(bool); ok {
			allowed = val
		}
	}

	result := &PolicyResult{
		Allowed: allowed,
		Rules:   []string{"opa:rego"},
	}

	if !allowed {
		result.Reason = "denied by OPA policy"
	} else {
		result.Reason = "allowed by OPA policy"
	}

	return result, nil
}
