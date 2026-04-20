package handler

import (
	"encoding/json"
	"net/http"

	"github.com/hexzedels/gosdlworkshop/internal/auth"
)

// HandleAdminSessions handles GET /api/admin/sessions — lists all active sessions.
func HandleAdminSessions(w http.ResponseWriter, r *http.Request) {
	sessions := auth.ListSessions()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}
