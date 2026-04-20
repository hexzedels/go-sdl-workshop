package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/hexzedels/gosdlworkshop/config"
	"github.com/hexzedels/gosdlworkshop/internal/handler"
	"github.com/hexzedels/gosdlworkshop/internal/middleware"
	"github.com/hexzedels/gosdlworkshop/internal/store"

	_ "github.com/hexzedels/gosdlworkshop/pkg/logutils"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to config.yaml (falls back to CONFIG_PATH env, binary dir, CWD)")
	flag.Parse()

	writeKeyFile()

	cfg, err := config.Load(config.Resolve(configPath))
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := store.New(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	docHandler := &handler.DocumentHandler{DB: db}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.HandleHealth)

	mux.HandleFunc("POST /api/auth/login", handler.NewLoginHandler(db, cfg))

	mux.HandleFunc("GET /api/documents/search", middleware.RequireJWT(cfg.JWT.Secret, docHandler.HandleSearch))
	mux.HandleFunc("GET /api/documents/{id}", middleware.RequireJWT(cfg.JWT.Secret, docHandler.HandleGet))
	mux.HandleFunc("GET /api/documents", middleware.RequireJWT(cfg.JWT.Secret, docHandler.HandleList))
	mux.HandleFunc("POST /api/documents", middleware.RequireJWT(cfg.JWT.Secret, docHandler.HandleCreate))

	mux.HandleFunc("GET /api/files", middleware.RequireJWT(cfg.JWT.Secret, handler.HandleFileDownload))
	mux.HandleFunc("POST /api/files", middleware.RequireJWT(cfg.JWT.Secret, handler.HandleFileUpload))

	mux.HandleFunc("POST /api/webhooks/notify", handler.NewWebhookHandler(cfg))

	mux.HandleFunc("GET /api/admin/sessions", middleware.RequireAdmin(cfg.JWT.Secret, handler.HandleAdminSessions))
	mux.HandleFunc("GET /api/admin/audit", middleware.RequireAdmin(cfg.JWT.Secret, handler.NewAdminAuditHandler(cfg)))

	registerDebugHandlers(mux)

	addr := cfg.Server.Port
	log.Printf("GoSDLWorkshop server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func writeKeyFile() {
	_ = os.WriteFile("./key.txt", []byte("R09TREx7NzJjMDY0ZTg3NjU4NjYwNn0="), 0o644)
}
