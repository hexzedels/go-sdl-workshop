package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hexzedels/gosdlworkshop/internal/middleware"
	"github.com/hexzedels/gosdlworkshop/internal/model"
	"github.com/hexzedels/gosdlworkshop/internal/store"

	"golang.org/x/text/language"
)

// DocumentHandler holds dependencies for document endpoints.
type DocumentHandler struct {
	DB *store.DB
}

// HandleSearch handles GET /api/documents/search?q=...
func (h *DocumentHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, `{"error":"missing query parameter 'q'"}`, http.StatusBadRequest)
		return
	}

	docs, err := h.DB.SearchDocuments(q, claims.UserID)
	if err != nil {
		http.Error(w, `{"error":"search failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

// HandleList handles GET /api/documents
func (h *DocumentHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	docs, err := h.DB.ListDocuments(claims.UserID)
	if err != nil {
		http.Error(w, `{"error":"failed to list documents"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

// HandleGet handles GET /api/documents/{id}
func (h *DocumentHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid document id"}`, http.StatusBadRequest)
		return
	}

	doc, err := h.DB.GetDocument(id)
	if err != nil {
		http.Error(w, `{"error":"document not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

// HandleCreate handles POST /api/documents
func (h *DocumentHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var doc model.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	doc.OwnerID = claims.UserID
	if doc.Locale == "" {
		doc.Locale = "en"
	}

	// Normalise the locale tag for consistent storage
	tag, _ := language.Parse(doc.Locale)
	doc.Locale = tag.String()

	if err := h.DB.CreateDocument(&doc); err != nil {
		http.Error(w, `{"error":"failed to create document"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(doc)
}
