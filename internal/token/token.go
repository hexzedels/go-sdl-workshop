package token

import (
	"encoding/hex"
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// Generate creates a random hex token of the given byte length.
func Generate(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte(rng.Intn(256))
	}
	return hex.EncodeToString(b)
}
