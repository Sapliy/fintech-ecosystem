package domain

import (
	"fmt"
)

type Action string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
	ActionAll    Action = "*"
)

type Resource string

const (
	ResourceZone        Resource = "zone"
	ResourceFlow        Resource = "flow"
	ResourceLedger      Resource = "ledger"
	ResourceTransaction Resource = "transaction"
	ResourceAPIKey      Resource = "api_key"
)

// PolicyEngine handles RBAC checks.
type PolicyEngine struct{}

func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{}
}

// Can check if a user with a specific role can perform an action on a resource.
func (e *PolicyEngine) Can(role string, action Action, resource Resource) bool {
	switch role {
	case RoleOwner, RoleAdmin:
		return true // Admins and Owners can do everything
	case RoleFinance:
		// Finance can read everything but only modify financial resources
		if action == ActionRead {
			return true
		}
		return resource == ResourceLedger || resource == ResourceTransaction
	case RoleDeveloper:
		// Developers can create and manage zones and flows, but not view sensitive financial data or modify owners
		if resource == ResourceZone || resource == ResourceFlow || resource == ResourceAPIKey {
			return true
		}
		// Deny reading sensitive financial data unless explicitly allowed
		if resource == ResourceLedger || resource == ResourceTransaction {
			return action == ActionRead
		}
	case RoleMember:
		// Basic members can only read
		return action == ActionRead
	}
	return false
}

func (e *PolicyEngine) ValidateAction(role string, action Action, resource Resource) error {
	if !e.Can(role, action, resource) {
		return fmt.Errorf("role %s is not authorized to %s resource %s", role, action, resource)
	}
	return nil
}
