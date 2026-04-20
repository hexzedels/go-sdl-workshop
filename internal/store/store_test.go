package store

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// TestSeedCredentials verifies that the bcrypt hashes in seed.sql accept the
// documented seed passwords.
func TestSeedCredentials(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("open in-memory store: %v", err)
	}
	defer db.Close()

	cases := []struct {
		username, password string
	}{
		{"admin", "admin123"},
		{"user", "password123"},
	}
	for _, tc := range cases {
		u, err := db.GetUserByUsername(tc.username)
		if err != nil {
			t.Fatalf("lookup %s: %v", tc.username, err)
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(tc.password)); err != nil {
			t.Fatalf("seed password %q does not match hash for %s: %v", tc.password, tc.username, err)
		}
	}
}
