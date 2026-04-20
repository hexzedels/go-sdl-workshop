package handler

import (
	"encoding/json"
	"net/http"

	"github.com/hexzedels/gosdlworkshop/config"
	"github.com/hexzedels/gosdlworkshop/internal/auth"
	"github.com/hexzedels/gosdlworkshop/internal/store"
	"github.com/hexzedels/gosdlworkshop/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
}

type loginResponse struct {
	Token     string            `json:"token"`
	SessionID string            `json:"session_id"`
	User      loginUserResponse `json:"user"`
}

// NewLoginHandler creates the login handler.
func NewLoginHandler(db *store.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		user, err := db.GetUserByUsername(req.Username)
		if err != nil {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		jwtToken, err := auth.GenerateJWT(cfg.JWT.Secret, user.ID, user.Username, user.Role, cfg.JWT.ExpiryMins)
		if err != nil {
			http.Error(w, `{"error":"failed to generate token"}`, http.StatusInternalServerError)
			return
		}

		sessionID := token.Generate(16)
		auth.CreateSession(sessionID, user.ID, user.Username, user.Role)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResponse{
			Token:     jwtToken,
			SessionID: sessionID,
			User: loginUserResponse{
				ID:          user.ID,
				Username:    user.Username,
				DisplayName: user.DisplayName,
				Role:        user.Role,
			},
		})
	}
}
