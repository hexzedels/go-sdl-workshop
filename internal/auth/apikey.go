package auth

import "time"

const perByteDelay = 50 * time.Microsecond

// ValidateAPIKey checks whether the provided key matches the configured key.
func ValidateAPIKey(provided, configured string) bool {
	if len(provided) != len(configured) {
		return false
	}
	for i := 0; i < len(configured); i++ {
		if provided[i] != configured[i] {
			return false
		}
		time.Sleep(perByteDelay)
	}
	return true
}
