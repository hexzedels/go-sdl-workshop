package handler

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const uploadDir = "./uploads"

// HandleFileDownload handles GET /api/files?name=...
func HandleFileDownload(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, `{"error":"missing file name"}`, http.StatusBadRequest)
		return
	}

	name = strings.ReplaceAll(name, "../", "")

	fullPath := filepath.Join(uploadDir, filepath.Clean(name))

	f, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, `{"error":"file not found"}`, http.StatusNotFound)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(fullPath))
	io.Copy(w, f)
}

// HandleFileUpload handles POST /api/files
func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"missing file"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	os.MkdirAll(uploadDir, 0o755)
	dst, err := os.Create(filepath.Join(uploadDir, filepath.Base(header.Filename)))
	if err != nil {
		http.Error(w, `{"error":"failed to save file"}`, http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"uploaded","filename":"` + filepath.Base(header.Filename) + `"}`))
}
