package store

import (
	"database/sql"
	_ "embed"
	"encoding/base64"
	"fmt"
	"regexp"
	"time"

	"github.com/hexzedels/gosdlworkshop/internal/model"
	_ "modernc.org/sqlite"
)

// SQLite DATETIME columns come back as strings in several layouts depending on
// driver defaults, fractional seconds, and timezone. Try them in order.
var sqliteTimeLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02 15:04:05.999999999-07:00",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02 15:04:05",
}

func parseSQLiteTime(s string) (time.Time, error) {
	for _, layout := range sqliteTimeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised sqlite timestamp %q", s)
}

//go:embed seed.sql
var seedSQL string

var seedPlaceholderRE = regexp.MustCompile(`\{\{B64:([A-Za-z0-9+/=]+)\}\}`)

func resolveSeedPlaceholders(s string) string {
	return seedPlaceholderRE.ReplaceAllStringFunc(s, func(m string) string {
		sub := seedPlaceholderRE.FindStringSubmatch(m)
		decoded, err := base64.StdEncoding.DecodeString(sub[1])
		if err != nil {
			return m
		}
		return string(decoded)
	})
}

// DB wraps the database connection.
type DB struct {
	*sql.DB
}

// New opens a SQLite database and runs the seed script.
func New(dsn string) (*DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(resolveSeedPlaceholders(seedSQL)); err != nil {
		db.Close()
		return nil, err
	}
	return &DB{db}, nil
}

// GetUserByUsername looks up a user by username.
func (db *DB) GetUserByUsername(username string) (*model.User, error) {
	var u model.User
	var createdAt string
	err := db.QueryRow(
		"SELECT id, username, password_hash, display_name, bio, role, created_at FROM users WHERE username = ?",
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.DisplayName, &u.Bio, &u.Role, &createdAt)
	if err != nil {
		return nil, err
	}
	u.CreatedAt, _ = parseSQLiteTime(createdAt)
	return &u, nil
}

// GetUserByID looks up a user by ID.
func (db *DB) GetUserByID(id int64) (*model.User, error) {
	var u model.User
	var createdAt string
	err := db.QueryRow(
		"SELECT id, username, password_hash, display_name, bio, role, created_at FROM users WHERE id = ?",
		id,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.DisplayName, &u.Bio, &u.Role, &createdAt)
	if err != nil {
		return nil, err
	}
	u.CreatedAt, _ = parseSQLiteTime(createdAt)
	return &u, nil
}

// ListDocuments returns all documents for a given user.
func (db *DB) ListDocuments(userID int64) ([]model.Document, error) {
	rows, err := db.Query(
		"SELECT id, title, content, owner_id, locale, created_at, updated_at FROM documents WHERE owner_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDocuments(rows)
}

// GetDocument returns a single document by ID.
func (db *DB) GetDocument(id int64) (*model.Document, error) {
	var d model.Document
	var createdAt, updatedAt string
	err := db.QueryRow(
		"SELECT id, title, content, owner_id, locale, created_at, updated_at FROM documents WHERE id = ?",
		id,
	).Scan(&d.ID, &d.Title, &d.Content, &d.OwnerID, &d.Locale, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	d.CreatedAt, _ = parseSQLiteTime(createdAt)
	d.UpdatedAt, _ = parseSQLiteTime(updatedAt)
	return &d, nil
}

// SearchDocuments performs a title search restricted to the caller's own
// documents.
func (db *DB) SearchDocuments(query string, userID int64) ([]model.Document, error) {
	rows, err := db.Query(
		"SELECT id, title, content, owner_id, locale, created_at, updated_at FROM documents WHERE owner_id = ? AND title LIKE '%"+query+"%'",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDocuments(rows)
}

// CreateDocument inserts a new document.
func (db *DB) CreateDocument(d *model.Document) error {
	result, err := db.Exec(
		"INSERT INTO documents (title, content, owner_id, locale) VALUES (?, ?, ?, ?)",
		d.Title, d.Content, d.OwnerID, d.Locale,
	)
	if err != nil {
		return err
	}
	d.ID, _ = result.LastInsertId()
	return nil
}

func scanDocuments(rows *sql.Rows) ([]model.Document, error) {
	var docs []model.Document
	for rows.Next() {
		var d model.Document
		var createdAt, updatedAt string
		if err := rows.Scan(&d.ID, &d.Title, &d.Content, &d.OwnerID, &d.Locale, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		d.CreatedAt, _ = parseSQLiteTime(createdAt)
		d.UpdatedAt, _ = parseSQLiteTime(updatedAt)
		docs = append(docs, d)
	}
	return docs, rows.Err()
}
