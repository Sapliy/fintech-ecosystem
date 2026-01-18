package main

import (
	"encoding/json"
	"microservices/internal/auth"
	"microservices/pkg/bcryptutil"
	"microservices/pkg/jsonutil"
	"microservices/pkg/jwtutil"
	"net/http"
)

// AuthHandler holds dependencies for authentication endpoints.
type AuthHandler struct {
	repo *auth.Repository
}

// RegisterRequest defines the payload for user registration.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest defines the payload for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse defines the successful response for login.
type LoginResponse struct {
	Token string `json:"token"`
}

// Register handles user account creation.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonutil.WriteErrorJSON(w, "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.WriteErrorJSON(w, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		jsonutil.WriteErrorJSON(w, "Email and password are required")
		return
	}

	// Use bcryptutil to hash the password securely
	b := &bcryptutil.BcryptUtilsImpl{}
	passwordHash, err := b.GenerateHash(req.Password)
	if err != nil {
		jsonutil.WriteErrorJSON(w, "Failed to hash password")
		return
	}

	user, err := h.repo.CreateUser(r.Context(), req.Email, passwordHash)
	if err != nil {
		// Check for duplicate key error (generic check for now)
		// In a real app, parse the error code for unique constraint violation
		jsonutil.WriteErrorJSON(w, "Failed to create user (email might be taken)")
		return
	}

	jsonutil.WriteJSON(w, http.StatusCreated, user)
}

// Login handles user authentication and JWT issuance.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonutil.WriteErrorJSON(w, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.WriteErrorJSON(w, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		jsonutil.WriteErrorJSON(w, "Email and password are required")
		return
	}

	// Fetch user by email
	user, err := h.repo.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		jsonutil.WriteErrorJSON(w, "Internal server error")
		return
	}
	if user == nil {
		jsonutil.WriteErrorJSON(w, "Invalid email or password") // Generic message for security
		return
	}

	// Verify password
	b := &bcryptutil.BcryptUtilsImpl{}
	match := b.CompareHash(req.Password, user.Password)
	if !match {
		jsonutil.WriteErrorJSON(w, "Invalid email or password")
		return
	}

	// Generate JWT
	token, err := jwtutil.GenerateToken(user.ID, user.Email)
	if err != nil {
		jsonutil.WriteErrorJSON(w, "Failed to generate token")
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, LoginResponse{Token: token})
}
