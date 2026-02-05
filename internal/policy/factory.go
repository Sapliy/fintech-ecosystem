package policy

import (
	"log"
	"os"
)

// NewEngine returns the configured policy engine.
// It tries to load from config/policies.rego (OPA), config/policies.json (JSON),
// or paths in POLICY_FILE_REGO / POLICY_FILE_JSON env vars.
// Falls back to HardcodedPolicyEngine if both load fails.
func NewEngine() PolicyEngine {
	// Check for OPA first
	regoPath := os.Getenv("POLICY_FILE_REGO")
	if regoPath == "" {
		locations := []string{"config/policies.rego", "../config/policies.rego", "../../config/policies.rego", "policies.rego"}
		for _, loc := range locations {
			if _, err := os.Stat(loc); err == nil {
				regoPath = loc
				break
			}
		}
	}

	if regoPath != "" {
		engine, err := NewOPAPolicyEngine(regoPath)
		if err == nil {
			log.Printf("Using OPA Policy Engine with config: %s", regoPath)
			return engine
		}
		log.Printf("Warning: Failed to load OPA policies from %s: %v. Trying JSON fallback.", regoPath, err)
	}

	// Fallback to JSON
	jsonPath := os.Getenv("POLICY_FILE_JSON")
	if jsonPath == "" {
		locations := []string{"config/policies.json", "../config/policies.json", "../../config/policies.json", "policies.json"}
		for _, loc := range locations {
			if _, err := os.Stat(loc); err == nil {
				jsonPath = loc
				break
			}
		}
	}

	if jsonPath != "" {
		engine, err := NewJSONPolicyEngine(jsonPath)
		if err == nil {
			log.Printf("Using JSON Policy Engine with config: %s", jsonPath)
			return engine
		}
		log.Printf("Warning: Failed to load JSON policies from %s: %v. Falling back to hardcoded policies.", jsonPath, err)
	}

	log.Println("Using Hardcoded Policy Engine")
	return NewHardcodedPolicyEngine()
}
