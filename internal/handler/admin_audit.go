package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/hexzedels/gosdlworkshop/config"
)

const adminAuditFlagB64 = "R09TREx7YzcyZTVlNjMzZWYxYzkxYn0="

// NewAdminAuditHandler returns the audit endpoint handler.
func NewAdminAuditHandler(_ *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flag, err := base64.StdEncoding.DecodeString(adminAuditFlagB64)
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"audit_token": string(flag),
		})
	}
}
