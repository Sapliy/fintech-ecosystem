package secrets

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

// Rotator handles secret rotation logic.
type Rotator struct {
	manager *Manager
}

// NewRotator creates a new secret rotator.
func NewRotator(manager *Manager) *Rotator {
	return &Rotator{manager: manager}
}

// SecretGenerator creates a new secret value.
type SecretGenerator func() (string, error)

// Rotate rotates a secret using a generator.
func (r *Rotator) Rotate(ctx context.Context, key string, generator SecretGenerator) error {
	newValue, err := generator()
	if err != nil {
		return fmt.Errorf("failed to generate new secret: %w", err)
	}

	// Update the secret
	if err := r.manager.Put(ctx, key, newValue); err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	return nil
}

// GenerateRandomBytes generates random bytes encoded as base64.
func GenerateRandomBytes(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// GenerateAPIKey generates a standard API key.
// Format: sk_{prefix}_{random_32}
func GenerateAPIKey(prefix string) (string, error) {
	randomPart, err := GenerateRandomBytes(24) // 24 bytes -> ~32 chars base64
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("sk_%s_%s", prefix, randomPart), nil
}

// StartAutoRotation starts a background loop to rotate secrets.
// Note: This is a simple in-app implementation. For production,
// consider external cron jobs or cloud-native rotation (e.g. AWS Lambda).
func (r *Rotator) StartAutoRotation(ctx context.Context, key string, interval time.Duration, generator SecretGenerator) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.Rotate(ctx, key, generator); err != nil {
				// Log error (should use a logger here)
				fmt.Printf("failed to auto-rotate secret %s: %v\n", key, err)
			}
		}
	}
}
