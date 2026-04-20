package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hexzedels/gosdlworkshop/config"
	"github.com/hexzedels/gosdlworkshop/internal/auth"
	"github.com/hexzedels/gosdlworkshop/internal/model"
)

// NewWebhookHandler creates the webhook notification handler.
func NewWebhookHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if !auth.ValidateAPIKey(apiKey, cfg.APIKey) {
			http.Error(w, `{"error":"invalid api key"}`, http.StatusUnauthorized)
			return
		}

		var payload model.WebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		if len(payload.URLs) == 0 {
			http.Error(w, `{"error":"no urls provided"}`, http.StatusBadRequest)
			return
		}

		// Fire-and-forget goroutine per URL — no limit, no timeout
		for _, url := range payload.URLs {
			go func(u string) {
				resp, err := http.Post(u, "application/json",
					nil)
				if err != nil {
					fmt.Printf("webhook delivery failed for %s: %v\n", u, err)
					return
				}
				resp.Body.Close()
				fmt.Printf("webhook delivered to %s: %d\n", u, resp.StatusCode)
			}(url)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "accepted",
			"message": fmt.Sprintf("delivering to %d URLs", len(payload.URLs)),
		})
	}
}
